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
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/go-openapi/runtime/middleware"

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

	mux.Post("/rulechain", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_rulechain")(addRuleChainEndpoint(svc)),
		decodeNewRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/rulechain", kithttp.NewServer(
		kitot.TraceServer(tracer, "rulechain_list")(listRuleChainEndpoint(svc)),
		decodeListRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Get("/rulechain/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "rulechain_info")(rulechainInfoEndpoint(svc)),
		decodeRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Put("/rulechain/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_rulechain")(updateRuleChainEndpoint(svc)),
		decodeUpdateRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Delete("/rulechain/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_rulechain")(deleteRuleChainEndpoint(svc)),
		decodeRuleChainRequest,
		encodeResponse,
		opts...,
	))

	mux.Put("/updateRulechainStatus/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_rulechain_status")(updateRuleChainStatusEndpoint(svc)),
		decodeUpdateRuleChainStatusRequest,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/version", pandas.Version("rulechain"))
	mux.Handle("/metrics", promhttp.Handler())

	return SetupMiddleware(mux)
}

func assetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:    rulechain.Asset,
		AssetDir: rulechain.AssetDir,
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

func decodeNewRuleChainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, ErrUnsupportedContentType
	}

	var ruleChain rulechain.RuleChain
	if err := json.NewDecoder(r.Body).Decode(&ruleChain); err != nil {
		return nil, errors.Wrap(rulechain.ErrMalformedEntity, err)
	}

	return updateRuleChainReq{
		token:     r.Header.Get("Authorization"),
		rulechain: ruleChain,
	}, nil
}

func decodeListRuleChainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := listRuleChainReq{
		token: r.Header.Get("Authorization"),
	}
	return req, nil
}

func decodeRuleChainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := RuleChainInfoRequest{
		token:       r.Header.Get("Authorization"),
		RuleChainID: r.Header.Get("RuleChainID"),
	}
	return req, nil
}

func decodeUpdateRuleChainStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := updateRuleChainStatusRequest{
		token:        r.Header.Get("Authorization"),
		RuleChainID:  r.Header.Get("RuleChainID"),
		updatestatus: r.Header.Get("updatestatus"),
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
		token:     r.Header.Get("Authorization"),
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
