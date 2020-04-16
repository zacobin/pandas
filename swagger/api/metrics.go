// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"time"

	swagger "github.com/cloustone/pandas/swagger"
	"github.com/go-kit/kit/metrics"
)

var _ swagger.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     swagger.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc swagger.Service, counter metrics.Counter, latency metrics.Histogram) swagger.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) RetrieveSwaggerConfigs(ctx context.Context, token string) (viewed swagger.Configs, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_configs").Add(1)
		ms.latency.With("method", "view_configs").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveSwaggerConfigs(ctx, token)
}

func (ms *metricsMiddleware) RetrieveDownstreamSwagger(ctx context.Context, token, module string) (viewed swagger.DownstreamSwagger, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_swagger").Add(1)
		ms.latency.With("method", "view_swagger").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveDownstreamSwagger(ctx, token, module)
}
