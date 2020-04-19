// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/cloustone/pandas/lbs"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveEntityOp               = "save_entity"
	saveEntitysOp              = "save_entity"
	updateEntityOp             = "update_entity"
	updateEntityKeyOp          = "update_entity_by_key"
	retrieveEntityByIDOp       = "retrieve_entity_by_id"
	retrieveEntityByKeyOp      = "retrieve_entity_by_key"
	retrieveAllEntitysOp       = "retrieve_all_entitys"
	retrieveEntitysByChannelOp = "retrieve_lbs_by_chan"
	removeEntityOp             = "remove_entity"
	retrieveEntityIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ lbs.EntityRepository = (*entityRepositoryMiddleware)(nil)
	_ lbs.EntityCache      = (*entityCacheMiddleware)(nil)
)

type entityRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   lbs.EntityRepository
}

// EntityRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func EntityRepositoryMiddleware(tracer opentracing.Tracer, repo lbs.EntityRepository) lbs.EntityRepository {
	return entityRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm entityRepositoryMiddleware) Save(ctx context.Context, ths ...lbs.Entity) ([]lbs.Entity, error) {
	span := createSpan(ctx, trm.tracer, saveEntitysOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm entityRepositoryMiddleware) Update(ctx context.Context, th lbs.Entity) error {
	span := createSpan(ctx, trm.tracer, updateEntityOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm entityRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (lbs.Entity, error) {
	span := createSpan(ctx, trm.tracer, retrieveEntityByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm entityRepositoryMiddleware) RetrieveByAttribute(ctx context.Context, channel, subtopic string) ([]lbs.Entity, error) {
	span := createSpan(ctx, trm.tracer, retrieveEntityByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByAttribute(ctx, channel, subtopic)
}

func (trm entityRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata lbs.Metadata) (lbs.EntitysPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllEntitysOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm entityRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeEntityOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type entityCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  lbs.EntityCache
}

// EntityCacheMiddleware tracks request and their latency, and adds spans
// to context.
func EntityCacheMiddleware(tracer opentracing.Tracer, cache lbs.EntityCache) lbs.EntityCache {
	return entityCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm entityCacheMiddleware) Save(ctx context.Context, entityKey string, entityID string) error {
	span := createSpan(ctx, tcm.tracer, saveEntityOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, entityKey, entityID)
}

func (tcm entityCacheMiddleware) ID(ctx context.Context, entityKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveEntityIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, entityKey)
}

func (tcm entityCacheMiddleware) Remove(ctx context.Context, entityID string) error {
	span := createSpan(ctx, tcm.tracer, removeEntityOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, entityID)
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
