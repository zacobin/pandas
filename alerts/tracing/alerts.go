// SPDX-License-Identifier: Apache-2.0

// Package tracing contains middlewares that will add spans
// to existing traces.
package tracing

import (
	"context"

	"github.com/cloustone/pandas/alerts"
	opentracing "github.com/opentracing/opentracing-go"
)

var _ alerts.AlertRepository = (*alertRepositoryMiddleware)(nil)

type alertRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   alerts.AlertRepository
}

// AlertRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func AlertRepositoryMiddleware(repo alerts.AlertRepository, tracer opentracing.Tracer) alerts.AlertRepository {
	return alertRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (arm alertRepositoryMiddleware) Save(ctx context.Context, alert ...alerts.Alert) ([]alerts.Alert, error) {
	span := createSpan(ctx, arm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.Save(ctx, alert...)
}

func (arm alertRepositoryMiddleware) Update(ctx context.Context, alert alerts.Alert) error {
	span := createSpan(ctx, arm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.Update(ctx, alert)
}

func (arm alertRepositoryMiddleware) Retrieve(ctx context.Context, owner, name string) (alerts.Alert, error) {
	span := createSpan(ctx, arm.tracer, retrieveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.Retrieve(ctx, owner, name)
}

func (arm alertRepositoryMiddleware) Revoke(ctx context.Context, owner, name string) error {
	span := createSpan(ctx, arm.tracer, revokeOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.Revoke(ctx, owner, name)
}

func (arm alertRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, meta Metadata) (alerts.AlertsPage, error) {
	span := createSpan(ctx, arm.tracer, listOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.List(ctx, owner, offset, limit, name, meta)
}
