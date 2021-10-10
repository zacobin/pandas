// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux"
	swagger "github.com/cloustone/pandas/swagger"
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
func MakeHandler(tracer opentracing.Tracer, svc swagger.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Get("/swaggers", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_swagger")(listSwaggerEndpoint(svc)),
		decodeListView,
		encodeResponse,
		opts...,
	))

	r.Get("/swaggers/:service", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_swagger")(viewSwaggerEndpoint(svc)),
		decodeView,
		encodeSwaggerResponse,
		opts...,
	))

	r.GetFunc("/version", pandas.Version("swagger"))
	r.Handle("/metrics", promhttp.Handler())
	//r.Handle("/swag", SetupMiddleware(promhttp.Handler()))
	// r.Handle("/swag", SSwagger(promhttp.Handler()))
	// r.Handle("/swag/docs", RedocUI(promhttp.Handler()))

	//return r
	return SetupMiddleware(r)
}

func assetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
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
			fmt.Println(r.URL.Path)
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

func decodeView(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewSwaggerReq{
		token:   r.Header.Get("Authorization"),
		service: bone.GetValue(r, "service"),
		httpreq: r,
	}

	return req, nil
}

func decodeListView(_ context.Context, r *http.Request) (interface{}, error) {
	req := listSwaggerReq{
		token: r.Header.Get("Authorization"),
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

func encodeSwaggerResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	resp := response.(viewSwaggerRes)

	return json.NewEncoder(w).Encode(resp.httpresp.Body)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case swagger.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case swagger.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case swagger.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case swagger.ErrConflict:
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
