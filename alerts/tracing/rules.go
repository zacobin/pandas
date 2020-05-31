// SPDX-License-Identifier: Apache-2.0

// Package tracing contains middlewares that will add spans
// to existing traces.
package tracing

import (
	"context"

	"github.com/cloustone/pandas/alerts"
	"github.com/cloustone/pandas/rules"
	opentracing "github.com/opentracing/opentracing-go"
)

var _ alerts.AlertRuleRepository = (*ruleRepositoryMiddleware)(nil)

type ruleRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   rules.AlertRuleRepository
}

// AlertRuleRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func AlertRuleRepositoryMiddleware(repo rules.AlertRuleRepository, tracer opentracing.Tracer) rules.AlertRuleRepository {
	return ruleRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (rrm ruleRepositoryMiddleware) Save(ctx context.Context, rule ...rules.AlertRule) ([]alerts.AlertRule, error) {
	span := createSpan(ctx, rrm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Save(ctx, rule...)
}

func (rrm ruleRepositoryMiddleware) Update(ctx context.Context, rule rules.AlertRule) error {
	span := createSpan(ctx, rrm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Update(ctx, rule)
}

func (rrm ruleRepositoryMiddleware) Retrieve(ctx context.Context, owner, name string) (rules.AlertRule, error) {
	span := createSpan(ctx, rrm.tracer, retrieveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Retrieve(ctx, owner, name)
}

func (rrm ruleRepositoryMiddleware) Revoke(ctx context.Context, owner, name string) error {
	span := createSpan(ctx, rrm.tracer, revokeOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, owner, span)

	return rrm.repo.Revoke(ctx, owner, name)
}

func (rrm ruleRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, meta Metadata) (alerts.AlertRulesPage, error) {
	span := createSpan(ctx, rrm.tracer, listOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.RetrieveAll(ctx, owner, offset, limit, name, meta)
}
