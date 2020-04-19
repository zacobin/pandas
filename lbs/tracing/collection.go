// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/cloustone/pandas/lbs"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveCollectionOp               = "save_collection"
	saveCollectionsOp              = "save_collection"
	updateCollectionOp             = "update_collection"
	updateCollectionKeyOp          = "update_collection_by_key"
	retrieveCollectionByIDOp       = "retrieve_collection_by_id"
	retrieveCollectionByKeyOp      = "retrieve_collection_by_key"
	retrieveAllCollectionsOp       = "retrieve_all_collections"
	retrieveCollectionsByChannelOp = "retrieve_lbs_by_chan"
	removeCollectionOp             = "remove_collection"
	retrieveCollectionIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ lbs.CollectionRepository = (*collectionRepositoryMiddleware)(nil)
	_ lbs.CollectionCache      = (*collectionCacheMiddleware)(nil)
)

type collectionRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   lbs.CollectionRepository
}

// CollectionRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func CollectionRepositoryMiddleware(tracer opentracing.Tracer, repo lbs.CollectionRepository) lbs.CollectionRepository {
	return collectionRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm collectionRepositoryMiddleware) Save(ctx context.Context, ths ...lbs.Collection) ([]lbs.Collection, error) {
	span := createSpan(ctx, trm.tracer, saveCollectionsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm collectionRepositoryMiddleware) Update(ctx context.Context, th lbs.Collection) error {
	span := createSpan(ctx, trm.tracer, updateCollectionOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm collectionRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (lbs.Collection, error) {
	span := createSpan(ctx, trm.tracer, retrieveCollectionByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm collectionRepositoryMiddleware) RetrieveByAttribute(ctx context.Context, channel, subtopic string) ([]lbs.Collection, error) {
	span := createSpan(ctx, trm.tracer, retrieveCollectionByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByAttribute(ctx, channel, subtopic)
}

func (trm collectionRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata lbs.Metadata) (lbs.CollectionsPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllCollectionsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm collectionRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeCollectionOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type collectionCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  lbs.CollectionCache
}

// CollectionCacheMiddleware tracks request and their latency, and adds spans
// to context.
func CollectionCacheMiddleware(tracer opentracing.Tracer, cache lbs.CollectionCache) lbs.CollectionCache {
	return collectionCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm collectionCacheMiddleware) Save(ctx context.Context, collectionKey string, collectionID string) error {
	span := createSpan(ctx, tcm.tracer, saveCollectionOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, collectionKey, collectionID)
}

func (tcm collectionCacheMiddleware) ID(ctx context.Context, collectionKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveCollectionIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, collectionKey)
}

func (tcm collectionCacheMiddleware) Remove(ctx context.Context, collectionID string) error {
	span := createSpan(ctx, tcm.tracer, removeCollectionOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, collectionID)
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
