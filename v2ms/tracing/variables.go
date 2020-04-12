// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/cloustone/pandas/v2ms"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveVariableOp               = "save_variable"
	saveVariablesOp              = "save_variable"
	updateVariableOp             = "update_variable"
	updateVariableKeyOp          = "update_variable_by_key"
	retrieveVariableByIDOp       = "retrieve_variable_by_id"
	retrieveVariableByKeyOp      = "retrieve_variable_by_key"
	retrieveAllVariablesOp       = "retrieve_all_variables"
	retrieveVariablesByChannelOp = "retrieve_v2ms_by_chan"
	removeVariableOp             = "remove_variable"
	retrieveVariableIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ v2ms.VariableRepository = (*variableRepositoryMiddleware)(nil)
	_ v2ms.VariableCache      = (*variableCacheMiddleware)(nil)
)

type variableRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   v2ms.VariableRepository
}

// VariableRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func VariableRepositoryMiddleware(tracer opentracing.Tracer, repo v2ms.VariableRepository) v2ms.VariableRepository {
	return variableRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm variableRepositoryMiddleware) Save(ctx context.Context, ths ...v2ms.Variable) ([]v2ms.Variable, error) {
	span := createSpan(ctx, trm.tracer, saveVariablesOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm variableRepositoryMiddleware) Update(ctx context.Context, th v2ms.Variable) error {
	span := createSpan(ctx, trm.tracer, updateVariableOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm variableRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (v2ms.Variable, error) {
	span := createSpan(ctx, trm.tracer, retrieveVariableByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm variableRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata v2ms.Metadata) (v2ms.VariablesPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllVariablesOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm variableRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeVariableOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type variableCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  v2ms.VariableCache
}

// VariableCacheMiddleware tracks request and their latency, and adds spans
// to context.
func VariableCacheMiddleware(tracer opentracing.Tracer, cache v2ms.VariableCache) v2ms.VariableCache {
	return variableCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm variableCacheMiddleware) Save(ctx context.Context, variableKey string, variableID string) error {
	span := createSpan(ctx, tcm.tracer, saveVariableOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, variableKey, variableID)
}

func (tcm variableCacheMiddleware) ID(ctx context.Context, variableKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveVariableIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, variableKey)
}

func (tcm variableCacheMiddleware) Remove(ctx context.Context, variableID string) error {
	span := createSpan(ctx, tcm.tracer, removeVariableOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, variableID)
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
