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
	"github.com/cloustone/pandas/pms"
	"github.com/cloustone/pandas/twins"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-openapi/runtime/middleware"
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
func MakeHandler(tracer opentracing.Tracer, svc pms.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/projects", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_project")(addProjectEndpoint(svc)),
		decodeProjectCreation,
		encodeResponse,
		opts...,
	))

	r.Put("/projects/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_project")(updateProjectEndpoint(svc)),
		decodeProjectUpdate,
		encodeResponse,
		opts...,
	))

	r.Get("/projects/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_project")(viewProjectEndpoint(svc)),
		decodeProject,
		encodeResponse,
		opts...,
	))

	r.Delete("/projects/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "remove_project")(removeProjectEndpoint(svc)),
		decodeProject,
		encodeResponse,
		opts...,
	))

	r.Get("/views", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_projects")(listProjectsEndpoint(svc)),
		decodeProjectList,
		encodeResponse,
		opts...,
	))

	// others handlers
	r.GetFunc("/version", pandas.Version("twins"))
	r.Handle("/metrics", promhttp.Handler())

	return SetupMiddleware(r)
}

func assetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:    pms.Asset,
		AssetDir: pms.AssetDir,
		Prefix:   "dist",
	}
}

//SSwagger s_swagger
func SSwagger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger" || r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger/", http.StatusFound)
			return
		}

		if strings.Index(r.URL.Path, "/swagger/") == 0 {
			http.StripPrefix("/swagger/", http.FileServer(assetFS())).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

//RedocUI docs to show redoc ui
func RedocUI(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := middleware.RedocOpts{
			Path:     "docs",
			SpecURL:  r.URL.Host + "/swagger/static/swagger/swagger.yaml",
			RedocURL: r.URL.Host + "/swagger/static/js/redoc.standalone.js",
			Title:    "swagger api",
		}

		middleware.Redoc(opts, handler).ServeHTTP(w, r)
		return
	})
}

//SetupMiddleware setupmiddleware
func SetupMiddleware(handler http.Handler) http.Handler {
	return SSwagger(
		RedocUI(handler),
	)
}

// Project handlers
func decodeProjectCreation(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := addProjectReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeProjectUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := updateProjectReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeProject(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewProjectReq{
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

	req := listProjectReq{
		token:    r.Header.Get("Authorization"),
		limit:    l,
		offset:   o,
		name:     n,
		metadata: m,
	}

	return req, nil
}

func decodeProjectList(_ context.Context, r *http.Request) (interface{}, error) {
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

	req := listProjectReq{
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
