// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"time"

	"github.com/cloustone/pandas/pms"
	"github.com/go-kit/kit/metrics"
)

var _ pms.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     pms.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc pms.Service, counter metrics.Counter, latency metrics.Histogram) pms.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) AddProject(ctx context.Context, token string, project pms.Project) (saved pms.Project, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_project").Add(1)
		ms.latency.With("method", "add_project").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddProject(ctx, token, project)
}

func (ms *metricsMiddleware) UpdateProject(ctx context.Context, token string, project pms.Project) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_project").Add(1)
		ms.latency.With("method", "update_project").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateProject(ctx, token, project)
}

func (ms *metricsMiddleware) ViewProject(ctx context.Context, token, id string) (viewed pms.Project, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_project").Add(1)
		ms.latency.With("method", "view_project").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewProject(ctx, token, id)
}

func (ms *metricsMiddleware) ListProjects(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata pms.Metadata) (tw pms.ProjectsPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_projects").Add(1)
		ms.latency.With("method", "list_projects").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListProjects(ctx, token, offset, limit, name, metadata)
}

func (ms *metricsMiddleware) RemoveProject(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_project").Add(1)
		ms.latency.With("method", "remove_project").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveProject(ctx, token, id)
}
