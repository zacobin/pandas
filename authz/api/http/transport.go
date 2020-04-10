// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cloustone/pandas/authz"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/errors"

	log "github.com/cloustone/pandas/pkg/logger"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const contentType = "application/json"

var (
	// ErrUnsupportedContentType indicates unacceptable or lack of Content-Type
	ErrUnsupportedContentType = errors.New("unsupported content type")
	errMissingRefererHeader   = errors.New("missing referer header")
	errInvalidToken           = errors.New("invalid token")
	errNoTokenSupplied        = errors.New("no token supplied")
	// ErrFailedDecode indicates failed to decode request body
	ErrFailedDecode = errors.New("failed to decode request body")
	logger          log.Logger
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc authz.Service, tracer opentracing.Tracer, l log.Logger) http.Handler {
	logger = l

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	mux := bone.New()

	mux.Get("/roles", kithttp.NewServer(
		kitot.TraceServer(tracer, "roles")(rolesInfoEndpoint(svc)),
		decodeRolesInfoRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/realms/:roleName", kithttp.NewServer(
		kitot.TraceServer(tracer, "role_info")(roleInfoEndpoint(svc)),
		decodeRoleInfoRequest,
		encodeResponse,
		opts...,
	))

	mux.Patch("/roles", kithttp.NewServer(
		kitot.TraceServer(tracer, "role_update")(updateRoleEndpoint(svc)),
		decodeUpdateRoleRequest,
		encodeResponse,
		opts...,
	))

	mux.Post("/authz", kithttp.NewServer(
		kitot.TraceServer(tracer, "authz")(authzEndpoint(svc)),
		decodeAuthzRequest,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/version", mainflux.Version("realms"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func decodeRolesInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}

	return genericRequest{
		token: r.Header.Get("Authorization"),
	}, nil
}

func decodeRoleInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}
	vals := bone.GetQuery(r, "rooleName")
	if len(vals) > 1 {
		return nil, errInvalidToken
	}

	if len(vals) == 0 {
		return "", errNoTokenSupplied
	}
	roleName := vals[0]
	req := viewRoleInfoRequest{
		token:    r.Header.Get("Authorization"),
		roleName: roleName,
	}
	return req, nil
}

func decodeUpdateRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}
	role := authz.Role{}
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode role: %s", err))
		return nil, err
	}

	req := updateRoleRequest{
		token: r.Header.Get("Authorization"),
		role:  role,
	}
	return req, nil
}

func decodeAuthzRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}

	req := authzRequest{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode authz request: %s", err))
		return nil, err
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
	case authz.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case authz.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case authz.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case authz.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case io.EOF, io.ErrUnexpectedEOF:
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
