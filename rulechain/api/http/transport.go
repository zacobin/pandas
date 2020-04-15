// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

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
	"github.com/cloustone/pandas/rulechain"

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
func MakeHandler(svc rulechain.Service, tracer opentracing.Tracer, l log.Logger) http.Handler {
	logger = l

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	mux := bone.New()

	mux.Put("/realms/:userID", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_rulechain")(addRuleChainEndpoint(svc)),
		decodeNewRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/rulechain/:userID", kithttp.NewServer(
		kitot.TraceServer(tracer, "rulechain_list")(listRuleChainEndpoint(svc)),
		decodeListRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/rulechain/:userID", kithttp.NewServer(
		kitot.TraceServer(tracer, "rulechain_info")(rulechainInfoEndpoint(svc)),
		decodeRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Put("/realms/:userID", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_rulechain")(updateRuleChainEndpoint(svc)),
		decodeUpdateRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Delete("/realms/:userID", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_rulechain")(deleteRuleChainEndpoint(svc)),
		decodeRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Put("/realms/:userID", kithttp.NewServer(
		kitot.TraceServer(tracer, "start_rulechain")(startRuleChainEndpoint(svc)),
		decodeRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Put("/realms/:userID", kithttp.NewServer(
		kitot.TraceServer(tracer, "stop_rulechain")(stopRuleChainEndpoint(svc)),
		decodeRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/version", pandas.Version("rulechain"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func decodeNewRuleChainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}

	var ruleChain rulechain.RuleChain
	if err := json.NewDecoder(r.Body).Decode(&ruleChain); err != nil {
		return nil, errors.Wrap(rulechain.ErrMalformedEntity, err)
	}

	return updateRuleChainReq{
		rulechain: ruleChain,
	}, nil
}

func decodeListRuleChainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := listRuleChainReq{
		UserID: r.Header.Get("UserID"),
	}
	return req, nil
}

func decodeRuleChainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := RuleChainRequestInfo{
		UserID:      r.Header.Get("UserID"),
		RuleChainID: r.Header.Get("RuleChainID"),
	}
	return req, nil
}

func decodeUpdateRuleChainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ruleChain := rulechain.RuleChain{}
	if err := json.NewDecoder(r.Body).Decode(&ruleChain); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode rulechain: %s", err))
		return nil, err
	}
	req := updateRuleChainReq{
		rulechain: ruleChain,
	}
	return req, nil
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
		case errors.Contains(errorVal, rulechain.ErrMalformedEntity):
			w.WriteHeader(http.StatusBadRequest)
			logger.Warn(fmt.Sprintf("Failed to decode rulechain credentials: %s", errorVal))
		case errors.Contains(errorVal, rulechain.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusForbidden)
		case errors.Contains(errorVal, rulechain.ErrConflict):
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
		case errors.Contains(errorVal, rulechain.ErrRuleChainNotFound):
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
