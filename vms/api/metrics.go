// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"time"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/vms"
	"github.com/go-kit/kit/metrics"
)

var _ vms.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     vms.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc vms.Service, counter metrics.Counter, latency metrics.Histogram) vms.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) AddView(ctx context.Context, token string, view vms.View) (saved vms.View, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_view").Add(1)
		ms.latency.With("method", "add_view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddView(ctx, token, view)
}

func (ms *metricsMiddleware) UpdateView(ctx context.Context, token string, view vms.View) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_view").Add(1)
		ms.latency.With("method", "update_view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateView(ctx, token, view)
}

func (ms *metricsMiddleware) ViewView(ctx context.Context, token, id string) (viewed vms.View, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_view").Add(1)
		ms.latency.With("method", "view_view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewView(ctx, token, id)
}

func (ms *metricsMiddleware) ListViews(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata vms.Metadata) (tw vms.ViewsPage, err error) {
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
func (ms *metricsMiddleware) AddVariable(ctx context.Context, token string, variable vms.Variable) (saved vms.Variable, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_variable").Add(1)
		ms.latency.With("method", "add_variable").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddVariable(ctx, token, variable)
}

func (ms *metricsMiddleware) UpdateVariable(ctx context.Context, token string, variable vms.Variable) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_variable").Add(1)
		ms.latency.With("method", "update_variable").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateVariable(ctx, token, variable)
}

func (ms *metricsMiddleware) ViewVariable(ctx context.Context, token, id string) (viewed vms.Variable, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_variable").Add(1)
		ms.latency.With("method", "view_variable").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewVariable(ctx, token, id)
}

func (ms *metricsMiddleware) ListVariables(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata vms.Metadata) (tw vms.VariablesPage, err error) {
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

func (ms *metricsMiddleware) SaveStates(msg *mainflux.Message) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "save_states").Add(1)
		ms.latency.With("method", "save_states").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.SaveStates(msg)
}

// Models

func (ms *metricsMiddleware) AddModel(ctx context.Context, token string, model vms.Model) (saved vms.Model, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_model").Add(1)
		ms.latency.With("method", "add_model").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddModel(ctx, token, model)
}

func (ms *metricsMiddleware) UpdateModel(ctx context.Context, token string, model vms.Model) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_model").Add(1)
		ms.latency.With("method", "update_model").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateModel(ctx, token, model)
}

func (ms *metricsMiddleware) ViewModel(ctx context.Context, token, id string) (viewed vms.Model, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_model").Add(1)
		ms.latency.With("method", "view_model").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewModel(ctx, token, id)
}

func (ms *metricsMiddleware) ListModels(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata vms.Metadata) (tw vms.ModelsPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_models").Add(1)
		ms.latency.With("method", "list_models").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListModels(ctx, token, offset, limit, name, metadata)
}

func (ms *metricsMiddleware) RemoveModel(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_model").Add(1)
		ms.latency.With("method", "remove_model").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveModel(ctx, token, id)
}
