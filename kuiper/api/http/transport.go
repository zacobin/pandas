package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/kuiper"
	"github.com/cloustone/pandas/kuiper/util"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/mainflux"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	contentType = "application/json"
	offset      = "offset"
	limit       = "limit"
	name        = "name"
	metadata    = "metadata"
	defOffset   = 0
	defLimit    = 10
)

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	errInvalidQueryParams     = errors.New("invalid query params")
)

var (
	dataDir        string
	logger         = util.Log
	startTimeStamp int64
)

type statementDescriptor struct {
	Sql string `json:"sql,omitempty"`
}

func decodeStatementDescriptor(reader io.ReadCloser) (statementDescriptor, error) {
	sd := statementDescriptor{}
	err := json.NewDecoder(reader).Decode(&sd)
	// Problems decoding
	if err != nil {
		return sd, fmt.Errorf("Error decoding the statement descriptor: %v", err)
	}
	return sd, nil
}

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(tracer opentracing.Tracer, svc kuiper.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Get("/streams", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_streams")(listStreamsEndpoint(svc)),
		decodeStreamListing,
		encodeResponse,
		opts...,
	))

	r.Post("/streams", kithttp.NewServer(
		kitot.TraceServer(tracer, "create_stream")(createStreamEndpoint(svc)),
		decodeStreamCreation,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/:name", kithttp.NewServer(
		kitot.TraceServer(tracer, "get_stream")(viewStreamEndpoint(svc)),
		decodeStreamView,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/:name", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_stream")(deleteStreamEndpoint(svc)),
		decodeStreamDeletion,
		encodeResponse,
		opts...,
	))

	r.Get("/rules", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_rules")(listRulesEndpoint(svc)),
		decodeRuleListing,
		encodeResponse,
		opts...,
	))

	r.Post("/rules", kithttp.NewServer(
		kitot.TraceServer(tracer, "create_rule")(createRuleEndpoint(svc)),
		decodeRuleCreation,
		encodeResponse,
		opts...,
	))

	r.Get("/rules/:name", kithttp.NewServer(
		kitot.TraceServer(tracer, "get_rule")(viewRuleEndpoint(svc)),
		decodeRuleView,
		encodeResponse,
		opts...,
	))

	r.Get("/rules/:name", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_rule")(deleteRuleEndpoint(svc)),
		decodeRuleDeletion,
		encodeResponse,
		opts...,
	))

	r.Get("/rules/:name/status", kithttp.NewServer(
		kitot.TraceServer(tracer, "rule_status")(viewRuleStatusEndpoint(svc)),
		decodeRuleStatus,
		encodeResponse,
		opts...,
	))

	r.Post("/rules/:name/:action", kithttp.NewServer(
		kitot.TraceServer(tracer, "rule_start")(startRuleEndpoint(svc)),
		decodeRuleControl,
		encodeResponse,
		opts...,
	))

	// Plugins
	r.Get("/plugins/sources", kithttp.NewServer(
		kitot.TraceServer(tracer, "plugins_source_list")(listPluginSourcesEndpoint(svc)),
		decodePluginSourcesListing,
		encodeResponse,
		opts...,
	))

	r.Get("/plugins/sources/:name", kithttp.NewServer(
		kitot.TraceServer(tracer, "plugins_source_view")(viewPluginSourceEndpoint(svc)),
		decodePluginSourceView,
		encodeResponse,
		opts...,
	))

	r.Get("/plugins/sinks", kithttp.NewServer(
		kitot.TraceServer(tracer, "plugins_sinks_list")(listPluginSinksEndpoint(svc)),
		decodePluginSourcesListing,
		encodeResponse,
		opts...,
	))

	r.Get("/plugins/sinks/:name", kithttp.NewServer(
		kitot.TraceServer(tracer, "plugins_sinks_view")(viewPluginSinkEndpoint(svc)),
		decodePluginSourceView,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", pandas.Version("kuiper"))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(mainflux.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case kuiper.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case kuiper.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case kuiper.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case kuiper.ErrConflict:
		w.WriteHeader(http.StatusUnprocessableEntity)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case errInvalidQueryParams:
		w.WriteHeader(http.StatusBadRequest)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	default:
		switch err.(type) {
		case *json.SyntaxError:
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func readUintQuery(r *http.Request, key string, def uint64) (uint64, error) {
	vals := bone.GetQuery(r, key)
	if len(vals) > 1 {
		return 0, errInvalidQueryParams
	}

	if len(vals) == 0 {
		return def, nil
	}

	strval := vals[0]
	val, err := strconv.ParseUint(strval, 10, 64)
	if err != nil {
		return 0, errInvalidQueryParams
	}

	return val, nil
}

func readStringQuery(r *http.Request, key string) (string, error) {
	vals := bone.GetQuery(r, key)
	if len(vals) > 1 {
		return "", errInvalidQueryParams
	}

	if len(vals) == 0 {
		return "", nil
	}

	return vals[0], nil
}

func readMetadataQuery(r *http.Request, key string) (map[string]interface{}, error) {
	vals := bone.GetQuery(r, key)
	if len(vals) > 1 {
		return nil, errInvalidQueryParams
	}

	if len(vals) == 0 {
		return nil, nil
	}

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(vals[0]), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func decodeStreamListing(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listResourcesReq{token: r.Header.Get("Authorization")}
	return req, nil

}

func decodeStreamCreation(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := createStreamReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Name == "" || req.Json == "" {
		return nil, kuiper.ErrMalformedEntity
	}

	return req, nil
}

func decodeStreamView(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeStreamDeletion(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeRuleListing(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	o, err := readUintQuery(r, offset, defOffset)
	if err != nil {
		return nil, err
	}

	l, err := readUintQuery(r, limit, defLimit)
	if err != nil {
		return nil, err
	}

	n, err := readStringQuery(r, name)
	if err != nil {
		return nil, err
	}

	m, err := readMetadataQuery(r, "metadata")
	if err != nil {
		return nil, err
	}

	req := listResourcesReq{
		token:    r.Header.Get("Authorization"),
		offset:   o,
		limit:    l,
		name:     n,
		metadata: m,
	}
	return req, nil
}

func decodeRuleCreation(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := createRuleReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Name == "" || req.SQL == "" {
		return nil, kuiper.ErrMalformedEntity
	}
	return req, nil
}

func decodeRuleView(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeRuleDeletion(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeRuleStatus(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeRuleControl(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := ruleControlReq{
		token:  r.Header.Get("Authorization"),
		id:     bone.GetValue(r, "id"),
		action: bone.GetValue(r, "action"),
	}

	return req, nil
}

// Plugin
func decodePluginSourcesListing(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	o, err := readUintQuery(r, offset, defOffset)
	if err != nil {
		return nil, err
	}

	l, err := readUintQuery(r, limit, defLimit)
	if err != nil {
		return nil, err
	}

	n, err := readStringQuery(r, name)
	if err != nil {
		return nil, err
	}

	m, err := readMetadataQuery(r, "metadata")
	if err != nil {
		return nil, err
	}

	req := listResourcesReq{
		token:    r.Header.Get("Authorization"),
		offset:   o,
		limit:    l,
		name:     n,
		metadata: m,
	}
	return req, nil
}

func decodePluginSourceView(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodePluginSinksListing(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	o, err := readUintQuery(r, offset, defOffset)
	if err != nil {
		return nil, err
	}

	l, err := readUintQuery(r, limit, defLimit)
	if err != nil {
		return nil, err
	}

	n, err := readStringQuery(r, name)
	if err != nil {
		return nil, err
	}

	m, err := readMetadataQuery(r, "metadata")
	if err != nil {
		return nil, err
	}

	req := listResourcesReq{
		token:    r.Header.Get("Authorization"),
		offset:   o,
		limit:    l,
		name:     n,
		metadata: m,
	}
	return req, nil
}

func decodePluginSinkView(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}
