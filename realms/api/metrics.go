// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"time"

	"github.com/cloustone/pandas/realms"
	"github.com/go-kit/kit/metrics"
)

var _ realms.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     realms.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc realms.Service, counter metrics.Counter, latency metrics.Histogram) realms.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) Register(ctx context.Context, realm realms.Realm) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "register").Add(1)
		ms.latency.With("method", "register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Register(ctx, realm)
}

func (ms *metricsMiddleware) RealmInfo(ctx context.Context, token, name string) (realms.Realm, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "info").Add(1)
		ms.latency.With("method", "info").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RealmInfo(ctx, token, name)
}

func (ms *metricsMiddleware) UpdateRealm(ctx context.Context, token string, realm realms.Realm) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_realm").Add(1)
		ms.latency.With("method", "update_realm").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateRealm(ctx, token, realm)
}

func (ms *metricsMiddleware) RevokeRealm(ctx context.Context, token string, name string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "revoke_realm").Add(1)
		ms.latency.With("method", "revoke_realm").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RevokeRealm(ctx, token, name)
}

func (ms *metricsMiddleware) ListRealms(ctx context.Context, token string) ([]realms.Realm, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_realm").Add(1)
		ms.latency.With("method", "list_realm").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListRealms(ctx, token)
}

func (ms *metricsMiddleware) Identify(ctx context.Context, token string, principal realms.Principal) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "identify").Add(1)
		ms.latency.With("method", "identify").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Identify(ctx, token, principal)
}
