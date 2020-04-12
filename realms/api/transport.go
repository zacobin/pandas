// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/errors"
	"github.com/cloustone/pandas/realms"

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
func MakeHandler(svc realms.Service, tracer opentracing.Tracer, l log.Logger) http.Handler {
	logger = l

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	mux := bone.New()

	mux.Post("/realms", kithttp.NewServer(
		kitot.TraceServer(tracer, "register")(registrationEndpoint(svc)),
		decodeNewRealmRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/realms", kithttp.NewServer(
		kitot.TraceServer(tracer, "realm_list")(listRealmEndpoint(svc)),
		decodeListRealmRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/realms/:realmName", kithttp.NewServer(
		kitot.TraceServer(tracer, "realm_info")(realmInfoEndpoint(svc)),
		decodeRealmRequest,
		encodeResponse,
		opts...,
	))

	mux.Put("/realms/:realmName", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_realm")(updateRealmEndpoint(svc)),
		decodeUpdateRealmRequest,
		encodeResponse,
		opts...,
	))

	mux.Delete("/realms/:realmName", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_realm")(deleteRealmEndpoint(svc)),
		decodeRealmRequest,
		encodeResponse,
		opts...,
	))
	mux.Post("/principals", kithttp.NewServer(
		kitot.TraceServer(tracer, "principal")(principalAuthEndpoint(svc)),
		decodePrincipalAuthRequest,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/version", pandas.Version("realms"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func decodeNewRealmRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}

	var realm realms.Realm
	if err := json.NewDecoder(r.Body).Decode(&realm); err != nil {
		return nil, errors.Wrap(realms.ErrMalformedEntity, err)
	}

	return createRealmReq{
		realm: realm,
		token: r.Header.Get("Authorization"),
	}, nil
}

func decodeListRealmRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewRealmInfoReq{
		token: r.Header.Get("Authorization"),
	}
	return req, nil
}

func decodeRealmRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := realmRequestInfo{
		token:     r.Header.Get("Authorization"),
		realmName: r.Header.Get("realmName"),
	}
	return req, nil
}

func decodeUpdateRealmRequest(_ context.Context, r *http.Request) (interface{}, error) {
	realm := realms.Realm{}
	if err := json.NewDecoder(r.Body).Decode(&realm); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode realm: %s", err))
		return nil, err
	}
	req := updateRealmReq{
		token:     r.Header.Get("Authorization"),
		realmName: r.Header.Get("realmName"),
		realm:     realm,
	}
	return req, nil
}

func decodePrincipalAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	principal := realms.Principal{}
	if err := json.NewDecoder(r.Body).Decode(&principal); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode principal: %s", err))
		return nil, err
	}
	req := principalAuthRequest{
		token:     r.Header.Get("Authorization"),
		principal: principal,
	}
	return req, nil
}

func decodeToken(_ context.Context, r *http.Request) (interface{}, error) {
	vals := bone.GetQuery(r, "token")
	if len(vals) > 1 {
		return "", errInvalidToken
	}

	if len(vals) == 0 {
		return "", errNoTokenSupplied
	}
	t := vals[0]
	return t, nil

}
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if ar, ok := response.(mainflux.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch errorVal := err.(type) {
	case errors.Error:
		w.Header().Set("Content-Type", contentType)
		switch {
		case errors.Contains(errorVal, realms.ErrMalformedEntity):
			w.WriteHeader(http.StatusBadRequest)
			logger.Warn(fmt.Sprintf("Failed to decode realm credentials: %s", errorVal))
		case errors.Contains(errorVal, realms.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusForbidden)
		case errors.Contains(errorVal, realms.ErrConflict):
			w.WriteHeader(http.StatusConflict)
		case errors.Contains(errorVal, ErrUnsupportedContentType):
			w.WriteHeader(http.StatusUnsupportedMediaType)
			logger.Warn("Invalid or missing content type.")
		case errors.Contains(errorVal, ErrFailedDecode):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, io.ErrUnexpectedEOF):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, io.EOF):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, realms.ErrRealmNotFound):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		if errorVal.Msg() != "" {
			if err := json.NewEncoder(w).Encode(errorRes{Err: errorVal.Msg()}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
