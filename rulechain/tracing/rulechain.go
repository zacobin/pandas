// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package tracing contains middlewares that will add spans
// to existing traces.
package tracing

import (
	"context"

	"github.com/cloustone/pandas/rulechain"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveOp                     = "save_op"
	retrieveByIDOp             = "retrieve_by_id"
	updatePassword             = "update_password"
	sendPasswordReset          = "send_reset_password"
	revokeRuleChain            = "revoke_rulechain"
	listRuleChain              = "list_rulechain"
	retrieveRuleChainIDByKeyOp = "retrieve_id_by_key"
)

var _ rulechain.RuleChainRepository = (*rulechainRepositoryMiddleware)(nil)

type rulechainRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   rulechain.RuleChainRepository
}

// RulechainRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func RulechainRepositoryMiddleware(repo rulechain.RuleChainRepository, tracer opentracing.Tracer) rulechain.RuleChainRepository {
	return rulechainRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (urm rulechainRepositoryMiddleware) Save(ctx context.Context, rulechain rulechain.RuleChain) error {
	span := createSpan(ctx, urm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Save(ctx, rulechain)
}

func (urm rulechainRepositoryMiddleware) Update(ctx context.Context, rulechain rulechain.RuleChain) error {
	span := createSpan(ctx, urm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Update(ctx, rulechain)
}

func (urm rulechainRepositoryMiddleware) Retrieve(ctx context.Context, UserID string, RuleChainID string) (rulechain.RuleChain, error) {
	span := createSpan(ctx, urm.tracer, retrieveByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Retrieve(ctx, UserID, RuleChainID)
}

func (urm rulechainRepositoryMiddleware) Revoke(ctx context.Context, UserID string, RuleChainID string) error {
	span := createSpan(ctx, urm.tracer, revokeRuleChain)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Revoke(ctx, UserID, RuleChainID)
}

func (urm rulechainRepositoryMiddleware) List(ctx context.Context, UserID string) ([]rulechain.RuleChain, error) {
	span := createSpan(ctx, urm.tracer, listRuleChain)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.List(ctx, UserID)
}

func createSpan(ctx context.Context, tracer opentracing.Tracer, opName string) opentracing.Span {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		return tracer.StartSpan(
			opName,
			opentracing.ChildOf(parentSpan.Context()),
		)
	}
	return tracer.StartSpan(opName)
}

type rulechainCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  rulechain.RuleChainCache
}

func RuleChainCacheMiddleware(tracer opentracing.Tracer, cache rulechain.RuleChainCache) rulechain.RuleChainCache {
	return rulechainCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (rcm rulechainCacheMiddleware) Save(ctx context.Context, projectKey string, projectID string) error {
	span := createSpan(ctx, rcm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	return rcm.cache.Save(ctx, projectKey, projectID)
}

func (rcm rulechainCacheMiddleware) ID(ctx context.Context, projectKey string) (string, error) {
	span := createSpan(ctx, rcm.tracer, retrieveRuleChainIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	return rcm.cache.ID(ctx, projectKey)
}

func (rcm rulechainCacheMiddleware) Remove(ctx context.Context, projectID string) error {
	span := createSpan(ctx, rcm.tracer, revokeRuleChain)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	return rcm.cache.Remove(ctx, projectID)
}
