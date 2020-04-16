// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"time"

	swagger_helper "github.com/cloustone/pandas/swagger-helper"
	"github.com/go-kit/kit/metrics"
)

var _ swagger_helper.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     swagger_helper.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc swagger_helper.Service, counter metrics.Counter, latency metrics.Histogram) swagger_helper.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) RetrieveDownstreamSwagger(ctx context.Context, token, module string) (viewed swagger_helper.DownstreamSwagger, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_swagger").Add(1)
		ms.latency.With("method", "view_swagger").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveDownstreamSwagger(ctx, token, module)
}
