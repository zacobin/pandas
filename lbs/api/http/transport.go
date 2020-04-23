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
	"github.com/cloustone/pandas/lbs"
	"github.com/cloustone/pandas/mainflux"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-zoo/bone"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const contentType = "application/json"

var errUnsupportedContentType = errors.New("unsupported content type")

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(tracer opentracing.Tracer, svc lbs.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Get("/collections", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_collections")(listCollectionsEndpoint(svc)),
		decodeListCollections,
		encodeResponse,
		opts...,
	))

	r.Post("/circlefences", kithttp.NewServer(
		kitot.TraceServer(tracer, "create_circlegeofence")(createCircleGeofenceEndpoint(svc)),
		decodeCreateCircleGeofence,
		encodeResponse,
		opts...,
	))

	r.Put("/circlefences", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_circlegeofence")(updateCircleGeofenceEndpoint(svc)),
		decodeUpdateCircleGeofence,
		encodeResponse,
		opts...,
	))

	r.Delete("/circlefence/:fenceID", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_circlegeofence")(deleteGeofenceEndpoint(svc)),
		decodeDeleteGeofence,
		encodeResponse,
		opts...,
	))

	r.Get("/fences", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_geofences")(listGeofencesEndpoint(svc)),
		decodeListGeofences,
		encodeResponse,
		opts...,
	))

	r.Post("/monitors/:fenceID", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_monitoredobject")(addMonitoredObjectEndpoint(svc)),
		decodeAddMonitoredObject,
		encodeResponse,
		opts...,
	))

	r.Delete("/monitors/:fenceID", kithttp.NewServer(
		kitot.TraceServer(tracer, "remove_monitoredobject")(removeMonitoredObjectEndpoint(svc)),
		decodeRemoveMonitoredObject,
		encodeResponse,
		opts...,
	))

	r.Get("/monitors/:fenceID", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_monitoredobject")(listMonitoredObjectsEndpoint(svc)),
		decodeListMonitoredObjects,
		encodeResponse,
		opts...,
	))

	r.Post("/polyfences", kithttp.NewServer(
		kitot.TraceServer(tracer, "create_polygeofence")(createPolyGeofenceEndpoint(svc)),
		decodeCreatePolyGeofence,
		encodeResponse,
		opts...,
	))

	r.Put("/polyfences", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_polygeofence")(updatePolyGeofenceEndpoint(svc)),
		decodeUpdatePolyGeofence,
		encodeResponse,
		opts...,
	))

	r.Get("/geofences", kithttp.NewServer(
		kitot.TraceServer(tracer, "get_fenceids")(getFenceIDsEndpoint(svc)),
		decodeGetFenceIDs,
		encodeResponse,
		opts...,
	))

	r.Get("/alarms/status", kithttp.NewServer(
		kitot.TraceServer(tracer, "query_status")(queryStatusEndpoint(svc)),
		decodeQueryStatus,
		encodeResponse,
		opts...,
	))

	r.Get("/historyalarms", kithttp.NewServer(
		kitot.TraceServer(tracer, "get_historyalarms")(getHistoryAlarmsEndpoint(svc)),
		decodeGetHistoryAlarms,
		encodeResponse,
		opts...,
	))

	r.Get("/historyalarms?batch", kithttp.NewServer(
		kitot.TraceServer(tracer, "batch_get_historyalarms")(batchGetHistoryAlarmsEndpoint(svc)),
		decodeBatchGetHistoryAlarms,
		encodeResponse,
		opts...,
	))

	r.Get("/staypoints", kithttp.NewServer(
		kitot.TraceServer(tracer, "get_staypoints")(getStayPointsEndpoint(svc)),
		decodeGetStayPoints,
		encodeResponse,
		opts...,
	))

	r.Post("/alarms", kithttp.NewServer(
		kitot.TraceServer(tracer, "notify_alarms")(notifyAlarmsEndpoint(svc)),
		decodeNotifyAlarms,
		encodeResponse,
		opts...,
	))

	// TODO
	r.Get("/lbs/userID/:fenceID", kithttp.NewServer(
		kitot.TraceServer(tracer, "get_fenceuserid")(getFenceUserIDEndpoint(svc)),
		decodeGetFenceUserID,
		encodeResponse,
		opts...,
	))

	r.Post("/entities", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_entity")(addEntityEndpoint(svc)),
		decodeAddEntity,
		encodeResponse,
		opts...,
	))

	r.Put("/entities", kithttp.NewServer(
		kitot.TraceServer(tracer, "update_entity")(updateEntityEndpoint(svc)),
		decodeUpdateEntity,
		encodeResponse,
		opts...,
	))

	r.Delete("/entities/:entitieName", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_entity")(deleteEntityEndpoint(svc)),
		decodeDeleteEntity,
		encodeResponse,
		opts...,
	))

	r.Get("/entities", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_entity")(listEntityEndpoint(svc)),
		decodeListEntity,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", pandas.Version("lbs"))
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/swagger", RedocUI(promhttp.Handler()))

	return r
}

//RedocUI docs to show redoc ui
func RedocUI(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := middleware.RedocOpts{
			Path:     "lbsdocs",
			SpecURL:  r.URL.Host + "/lbs/swagger.yaml",
			RedocURL: r.URL.Host + "/swagger/static/js/redoc.standalone.js",
			Title:    "lbs api",
		}
		middleware.Redoc(opts, handler).ServeHTTP(w, r)
		return
	})
}

func decodeListCollections(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeCreateCircleGeofence(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateCircleGeofence(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeDeleteGeofence(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeListGeofences(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeAddMonitoredObject(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeRemoveMonitoredObject(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeListMonitoredObjects(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeCreatePolyGeofence(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdatePolyGeofence(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetFenceIDs(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeQueryStatus(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetHistoryAlarms(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeBatchGetHistoryAlarms(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetStayPoints(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeNotifyAlarms(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetFenceUserID(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeAddEntity(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateEntity(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeDeleteEntity(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeListEntity(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}
	req := listCollectionsReq{
		token: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	case lbs.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case lbs.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case lbs.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case lbs.ErrConflict:
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
