// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package httpapi

import (
	"context"
	"time"

	"github.com/cloustone/pandas/authz"
	"github.com/go-kit/kit/metrics"
)

var _ authz.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     authz.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc authz.Service, counter metrics.Counter, latency metrics.Histogram) authz.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) RetrieveRole(ctx context.Context, token, roleName string) (authz.Role, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve").Add(1)
		ms.latency.With("method", "retrieve").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveRole(ctx, token, roleName)
}

func (ms *metricsMiddleware) ListRoles(ctx context.Context, token string) ([]authz.Role, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list").Add(1)
		ms.latency.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListRoles(ctx, token)
}

func (ms *metricsMiddleware) UpdateRole(ctx context.Context, token string, role authz.Role) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update").Add(1)
		ms.latency.With("method", "update").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateRole(ctx, token, role)
}

func (ms *metricsMiddleware) Authorize(ctx context.Context, token string, roleName string, subject authz.Subject) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "authorize").Add(1)
		ms.latency.With("method", "authorize").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Authorize(ctx, token, roleName, subject)
}
