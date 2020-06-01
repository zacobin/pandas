// SPDX-License-Identifier: Apache-2.0

// Package tracing contains middlewares that will add spans
// to existing traces.
package tracing

import (
	"context"

	"github.com/cloustone/pandas/alerts"
	opentracing "github.com/opentracing/opentracing-go"
)

var _ alerts.AlertRuleRepository = (*ruleRepositoryMiddleware)(nil)

type ruleRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   alerts.AlertRuleRepository
}

// AlertRuleRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func AlertRuleRepositoryMiddleware(repo alerts.AlertRuleRepository, tracer opentracing.Tracer) alerts.AlertRuleRepository {
	return ruleRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (rrm ruleRepositoryMiddleware) Save(ctx context.Context, rule alerts.AlertRule) (alerts.AlertRule, error) {
	span := createSpan(ctx, rrm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Save(ctx, rule)
}

func (rrm ruleRepositoryMiddleware) Update(ctx context.Context, rule alerts.AlertRule) error {
	span := createSpan(ctx, rrm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Update(ctx, rule)
}

func (rrm ruleRepositoryMiddleware) Retrieve(ctx context.Context, owner, name string) (alerts.AlertRule, error) {
	span := createSpan(ctx, rrm.tracer, retrieveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Retrieve(ctx, owner, name)
}

func (rrm ruleRepositoryMiddleware) Revoke(ctx context.Context, owner, name string) error {
	span := createSpan(ctx, rrm.tracer, revokeOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Revoke(ctx, owner, name)
}

func (rrm ruleRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, meta alerts.Metadata) (alerts.AlertRulesPage, error) {
	span := createSpan(ctx, rrm.tracer, listOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.RetrieveAll(ctx, owner, offset, limit, name, meta)
}
