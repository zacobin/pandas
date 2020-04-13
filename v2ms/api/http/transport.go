// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/twins"
	"github.com/cloustone/pandas/v2ms"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	contentType = "application/json"

	offset   = "offset"
	limit    = "limit"
	name     = "name"
	metadata = "metadata"

	defLimit  = 10
	defOffset = 0
)

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	errInvalidQueryParams     = errors.New("invalid query params")
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(tracer opentracing.Tracer, svc v2ms.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	// View endpoints
	r.Post("/views", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_view")(addViewEndpoint(svc)),
		decodeViewCreation,
		encodeResponse,
		opts...,
	))

	r.Put("/views/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_view")(updateViewEndpoint(svc)),
		decodeViewUpdate,
		encodeResponse,
		opts...,
	))

	r.Get("/views/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_view")(viewViewEndpoint(svc)),
		decodeView,
		encodeResponse,
		opts...,
	))

	r.Delete("/views/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "remove_view")(removeViewEndpoint(svc)),
		decodeView,
		encodeResponse,
		opts...,
	))

	r.Get("/views", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_twins")(listViewsEndpoint(svc)),
		decodeViewList,
		encodeResponse,
		opts...,
	))

	// Variables
	r.Post("/vars", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_var")(addVariableEndpoint(svc)),
		decodeVariableCreation,
		encodeResponse,
		opts...,
	))

	r.Put("/vars/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_var")(updateVariableEndpoint(svc)),
		decodeVariableUpdate,
		encodeResponse,
		opts...,
	))

	r.Get("/vars/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_var")(viewVariableEndpoint(svc)),
		decodeVariable,
		encodeResponse,
		opts...,
	))

	r.Delete("/vars/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "remove_var")(removeVariableEndpoint(svc)),
		decodeVariable,
		encodeResponse,
		opts...,
	))

	r.Get("/vars", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_vars")(listVariablesEndpoint(svc)),
		decodeVariableList,
		encodeResponse,
		opts...,
	))

	// Models
	r.Post("/models", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_models")(addModelEndpoint(svc)),
		decodeModelCreation,
		encodeResponse,
		opts...,
	))

	r.Put("/models/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_models")(updateModelEndpoint(svc)),
		decodeModelUpdate,
		encodeResponse,
		opts...,
	))

	r.Get("/models/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_model")(viewModelEndpoint(svc)),
		decodeModel,
		encodeResponse,
		opts...,
	))

	r.Delete("/models/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "remove_models")(removeModelEndpoint(svc)),
		decodeModel,
		encodeResponse,
		opts...,
	))

	r.Get("/models", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_models")(listModelsEndpoint(svc)),
		decodeModelList,
		encodeResponse,
		opts...,
	))

	// others handlers
	r.GetFunc("/version", pandas.Version("twins"))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

// View handlers
func decodeViewCreation(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := addViewReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeViewUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := updateViewReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeView(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewViewReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeViewList(_ context.Context, r *http.Request) (interface{}, error) {
	l, err := readUintQuery(r, limit, defLimit)
	if err != nil {
		return nil, err
	}

	o, err := readUintQuery(r, offset, defOffset)
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

	req := listViewReq{
		token:    r.Header.Get("Authorization"),
		limit:    l,
		offset:   o,
		name:     n,
		metadata: m,
	}

	return req, nil
}

//  Variables
func decodeVariableCreation(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := addVariableReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeVariableUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := updateViewReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeVariable(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewVariableReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}

	return req, nil
}

func decodeVariableList(_ context.Context, r *http.Request) (interface{}, error) {
	l, err := readUintQuery(r, limit, defLimit)
	if err != nil {
		return nil, err
	}

	o, err := readUintQuery(r, offset, defOffset)
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

	req := listVariableReq{
		token:    r.Header.Get("Authorization"),
		limit:    l,
		offset:   o,
		name:     n,
		metadata: m,
	}

	return req, nil
}

// Models

func decodeModelCreation(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := addModelReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeModelUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := updateModelReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeModel(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewVariableReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}

	return req, nil
}

func decodeModelList(_ context.Context, r *http.Request) (interface{}, error) {
	l, err := readUintQuery(r, limit, defLimit)
	if err != nil {
		return nil, err
	}

	o, err := readUintQuery(r, offset, defOffset)
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

	req := listModelReq{
		token:    r.Header.Get("Authorization"),
		limit:    l,
		offset:   o,
		name:     n,
		metadata: m,
	}

	return req, nil
}

// common

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
	case twins.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case twins.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case twins.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case twins.ErrConflict:
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
