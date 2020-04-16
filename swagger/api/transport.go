// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux"
	swagger "github.com/cloustone/pandas/swagger"
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

	return r
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
