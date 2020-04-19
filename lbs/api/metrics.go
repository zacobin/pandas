// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"time"

	"github.com/cloustone/pandas/lbs"
	"github.com/go-kit/kit/metrics"
)

var _ lbs.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     lbs.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc lbs.Service, counter metrics.Counter, latency metrics.Histogram) lbs.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
func (ms *metricsMiddleware) ListCollections(ctx context.Context, token string) (products []string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_collections").Add(1)
		ms.latency.With("method", "list_collections").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListCollections(ctx, token)
}

// Geofence
func (ms *metricsMiddleware) CreateCircleGeofence(ctx context.Context, token string, fence *lbs.CircleGeofence) (fenceID string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_circle_geofence").Add(1)
		ms.latency.With("method", "create_circle_geofence").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateCircleGeofence(ctx, token, fence)
}

func (ms *metricsMiddleware) UpdateCircleGeofence(ctx context.Context, token string, fence *lbs.CircleGeofence) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_circle_geofence").Add(1)
		ms.latency.With("method", "update_circle_geofence").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateCircleGeofence(ctx, token, fence)
}

func (ms *metricsMiddleware) DeleteGeofence(ctx context.Context, token string, fenceIDs []string, objects []string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "delete_geofence").Add(1)
		ms.latency.With("method", "delete_geofence").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.DeleteGeofence(ctx, token, fenceIDs, objects)
}

func (ms *metricsMiddleware) ListGeofences(ctx context.Context, token string, fenceIDs []string, objects []string) (fenceList []*lbs.Geofence, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_geofence").Add(1)
		ms.latency.With("method", "list_geofence").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListGeofences(ctx, token, fenceIDs, objects)
}

func (ms *metricsMiddleware) AddMonitoredObject(ctx context.Context, token string, fenceID string, objects []string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_monitored_object").Add(1)
		ms.latency.With("method", "add_monitored_object").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddMonitoredObject(ctx, token, fenceID, objects)
}

func (ms *metricsMiddleware) RemoveMonitoredObject(ctx context.Context, token string, fenceID string, objects []string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_monitored_object").Add(1)
		ms.latency.With("method", "remove_monitored_object").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveMonitoredObject(ctx, token, fenceID, objects)
}

func (ms *metricsMiddleware) ListMonitoredObjects(ctx context.Context, token string, fenceID string, pageIndex int, pageSize int) (total int, objects []string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_monitored_object").Add(1)
		ms.latency.With("method", "list_monitored_object").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.ListMonitoredObjects(ctx, token, fenceID, pageIndex, pageSize)
}

func (ms *metricsMiddleware) CreatePolyGeofence(ctx context.Context, token string, fence *lbs.PolyGeofence) (fenceID string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_poly_geofence").Add(1)
		ms.latency.With("method", "create_poly_geofence").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreatePolyGeofence(ctx, token, fence)
}

func (ms *metricsMiddleware) UpdatePolyGeofence(ctx context.Context, token string, fence *lbs.PolyGeofence) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_poly_geofence").Add(1)
		ms.latency.With("method", "update_poly_geofence").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdatePolyGeofence(ctx, token, fence)
}

func (ms *metricsMiddleware) GetFenceIDs(ctx context.Context, token string) (fenceIDs []string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_fenceids").Add(1)
		ms.latency.With("method", "get_fenceids").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetFenceIDs(ctx, token)
}

// Alarm
func (ms *metricsMiddleware) QueryStatus(ctx context.Context, token string, monitoredPerson string, fenceIDs []string) (status *lbs.QueryStatus, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "query_status").Add(1)
		ms.latency.With("method", "query_status").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.QueryStatus(ctx, token, monitoredPerson, fenceIDs)
}

func (ms *metricsMiddleware) GetHistoryAlarms(ctx context.Context, token string, monitoredPerson string, fenceIDs []string) (alarms *lbs.HistoryAlarms, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_history_alarms").Add(1)
		ms.latency.With("method", "get_history_alarms").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetHistoryAlarms(ctx, token, monitoredPerson, fenceIDs)
}

func (ms *metricsMiddleware) BatchGetHistoryAlarms(ctx context.Context, token string, input *lbs.BatchGetHistoryAlarmsRequest) (alarms *lbs.BatchHistoryAlarmsResp, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "batchget_history_alarms").Add(1)
		ms.latency.With("method", "batchget_history_alarms").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.BatchGetHistoryAlarms(ctx, token, input)
}

func (ms *metricsMiddleware) GetStayPoints(ctx context.Context, token string, input *lbs.GetStayPointsRequest) (points *lbs.StayPoints, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_stay_points").Add(1)
		ms.latency.With("method", "get_stay_points").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetStayPoints(ctx, token, input)
}

// NotifyAlarms is used by apiserver to provide asynchrous notication
func (ms *metricsMiddleware) NotifyAlarms(ctx context.Context, token string, content []byte) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "notify_alarm").Add(1)
		ms.latency.With("method", "notify_alarm").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.NotifyAlarms(ctx, token, content)
}

func (ms *metricsMiddleware) GetFenceUserID(ctx context.Context, token string, fenceID string) (userID string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_fence_userid").Add(1)
		ms.latency.With("method", "get_fence_userid").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetFenceUserID(ctx, token, fenceID)
}

//Entity
func (ms *metricsMiddleware) AddEntity(ctx context.Context, token string, entityName string, entityDesc string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_entity").Add(1)
		ms.latency.With("method", "add_entity").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddEntity(ctx, token, entityName, entityDesc)
}

func (ms *metricsMiddleware) UpdateEntity(ctx context.Context, token string, entityName string, entityDesc string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_entity").Add(1)
		ms.latency.With("method", "update_entity").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateEntity(ctx, token, entityName, entityDesc)
}

func (ms *metricsMiddleware) DeleteEntity(ctx context.Context, token string, entityName string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "delete_entity").Add(1)
		ms.latency.With("method", "delete_entity").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.DeleteEntity(ctx, token, entityName)
}

func (ms *metricsMiddleware) ListEntity(ctx context.Context, token string, coordTypeOutput string, pageIndex int, pageSize int) (total int, info []*lbs.EntityInfo, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_entity").Add(1)
		ms.latency.With("method", "list_entity").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListEntity(ctx, token, coordTypeOutput, pageIndex, pageSize)
}
