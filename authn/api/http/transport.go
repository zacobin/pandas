// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/authn"
	"github.com/cloustone/pandas/mainflux"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const contentType = "application/json"

var errUnsupportedContentType = errors.New("unsupported content type")

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc authn.Service, tracer opentracing.Tracer) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	mux := bone.New()

	mux.Post("/keys", kithttp.NewServer(
		kitot.TraceServer(tracer, "issue")(issueEndpoint(svc)),
		decodeIssue,
		encodeResponse,
		opts...,
	))

	mux.Get("/keys/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "retrieve")(retrieveEndpoint(svc)),
		decodeKeyReq,
		encodeResponse,
		opts...,
	))

	mux.Delete("/keys/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "revoke")(revokeEndpoint(svc)),
		decodeKeyReq,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/version", pandas.Version("auth"))
	mux.Handle("/metrics", promhttp.Handler())

	return SetupMiddleware(mux)
}

func assetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:    authn.Asset,
		AssetDir: authn.AssetDir,
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

func decodeIssue(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := issueKeyReq{
		issuer: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeKeyReq(_ context.Context, r *http.Request) (interface{}, error) {
	req := keyReq{
		issuer: r.Header.Get("Authorization"),
		id:     bone.GetValue(r, "id"),
	}
	return req, nil
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
	case authn.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case authn.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case authn.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case authn.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case io.EOF, io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
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
