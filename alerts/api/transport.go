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
	"github.com/cloustone/pandas/alerts"
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
func MakeHandler(svc alerts.Service, tracer opentracing.Tracer, l log.Logger) http.Handler {
	logger = l

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	mux := bone.New()

	mux.Post("/alerts", kithttp.NewServer(
		kitot.TraceServer(tracer, "create")(createAlertEndpoint(svc)),
		decodeNewAlertRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/alerts", kithttp.NewServer(
		kitot.TraceServer(tracer, "alerts_list")(listAlertsEndpoint(svc)),
		decodeListAlertsRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/alerts/:alertName", kithttp.NewServer(
		kitot.TraceServer(tracer, "realm_info")(alertInfoEndpoint(svc)),
		decodeAlertRequest,
		encodeResponse,
		opts...,
	))

	mux.Put("/alerts", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_alert")(updateAlertEndpoint(svc)),
		decodeNewAlertRequest,
		encodeResponse,
		opts...,
	))

	mux.Delete("/alerts/:alertName", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_alert")(deleteAlertEndpoint(svc)),
		decodeAlertRequest,
		encodeResponse,
		opts...,
	))

	mux.Post("/alertrules", kithttp.NewServer(
		kitot.TraceServer(tracer, "create")(createAlertRuleEndpoint(svc)),
		decodeNewAlertRuleRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/alertrules", kithttp.NewServer(
		kitot.TraceServer(tracer, "alertrule_list")(listAlertRulesEndpoint(svc)),
		decodeListAlertRulesRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/alertrules/:alertRuleName", kithttp.NewServer(
		kitot.TraceServer(tracer, "alertrule_info")(alertRuleInfoEndpoint(svc)),
		decodeAlertRequest,
		encodeResponse,
		opts...,
	))

	mux.Put("/alertrules", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_alertrule")(updateAlertRuleEndpoint(svc)),
		decodeNewAlertRuleRequest,
		encodeResponse,
		opts...,
	))

	mux.Delete("/alertrules/:alertRuleName", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_alertrule")(deleteAlertRuleEndpoint(svc)),
		decodeAlertRequest,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/version", pandas.Version("alerts"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func decodeNewAlertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}
	return nil, nil
}

func decodeListAlertsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}
	return nil, nil
}

func decodeAlertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}
	return nil, nil
}

func decodeNewAlertRuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}
	return nil, nil
}

func decodeListAlertRulesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}
	return nil, nil
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
		case errors.Contains(errorVal, alerts.ErrMalformedEntity):
			w.WriteHeader(http.StatusBadRequest)
			logger.Warn(fmt.Sprintf("Failed to decode realm credentials: %s", errorVal))
		case errors.Contains(errorVal, alerts.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusForbidden)
		case errors.Contains(errorVal, alerts.ErrConflict):
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
		case errors.Contains(errorVal, alerts.ErrAlertNotFound):
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
