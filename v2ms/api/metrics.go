// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"time"

	"github.com/cloustone/pandas/v2ms"
	"github.com/go-kit/kit/metrics"
)

var _ v2ms.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     v2ms.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc v2ms.Service, counter metrics.Counter, latency metrics.Histogram) v2ms.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) AddView(ctx context.Context, token string, view v2ms.View) (saved v2ms.View, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_view").Add(1)
		ms.latency.With("method", "add_view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddView(ctx, token, view)
}

func (ms *metricsMiddleware) UpdateView(ctx context.Context, token string, view v2ms.View) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_view").Add(1)
		ms.latency.With("method", "update_view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateView(ctx, token, view)
}

func (ms *metricsMiddleware) ViewView(ctx context.Context, token, id string) (viewed v2ms.View, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_view").Add(1)
		ms.latency.With("method", "view_view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewView(ctx, token, id)
}

func (ms *metricsMiddleware) ListViews(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata v2ms.Metadata) (tw v2ms.ViewsPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_views").Add(1)
		ms.latency.With("method", "list_views").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListViews(ctx, token, offset, limit, name, metadata)
}

func (ms *metricsMiddleware) RemoveView(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_view").Add(1)
		ms.latency.With("method", "remove_view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveView(ctx, token, id)
}

// Variable
func (ms *metricsMiddleware) AddVariable(ctx context.Context, token string, variable v2ms.Variable) (saved v2ms.Variable, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_variable").Add(1)
		ms.latency.With("method", "add_variable").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddVariable(ctx, token, variable)
}

func (ms *metricsMiddleware) UpdateVariable(ctx context.Context, token string, variable v2ms.Variable) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_variable").Add(1)
		ms.latency.With("method", "update_variable").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateVariable(ctx, token, variable)
}

func (ms *metricsMiddleware) ViewVariable(ctx context.Context, token, id string) (viewed v2ms.Variable, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_variable").Add(1)
		ms.latency.With("method", "view_variable").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewVariable(ctx, token, id)
}

func (ms *metricsMiddleware) ListVariables(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata v2ms.Metadata) (tw v2ms.VariablesPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_variables").Add(1)
		ms.latency.With("method", "list_variables").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListVariables(ctx, token, offset, limit, name, metadata)
}

func (ms *metricsMiddleware) RemoveVariable(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_variable").Add(1)
		ms.latency.With("method", "remove_variable").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveVariable(ctx, token, id)
}
