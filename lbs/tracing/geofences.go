// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/cloustone/pandas/lbs"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveGeofenceOp               = "save_geofence"
	saveGeofencesOp              = "save_geofence"
	updateGeofenceOp             = "update_geofence"
	updateGeofenceKeyOp          = "update_geofence_by_key"
	retrieveGeofenceByIDOp       = "retrieve_geofence_by_id"
	retrieveGeofenceByKeyOp      = "retrieve_geofence_by_key"
	retrieveAllGeofencesOp       = "retrieve_all_geofences"
	retrieveGeofencesByChannelOp = "retrieve_lbs_by_chan"
	removeGeofenceOp             = "remove_geofence"
	retrieveGeofenceIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ lbs.GeofenceRepository = (*geofenceRepositoryMiddleware)(nil)
	_ lbs.GeofenceCache      = (*geofenceCacheMiddleware)(nil)
)

type geofenceRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   lbs.GeofenceRepository
}

// GeofenceRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func GeofenceRepositoryMiddleware(tracer opentracing.Tracer, repo lbs.GeofenceRepository) lbs.GeofenceRepository {
	return geofenceRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm geofenceRepositoryMiddleware) Save(ctx context.Context, ths ...lbs.Geofence) ([]lbs.GeofenceRecord, error) {
	span := createSpan(ctx, trm.tracer, saveGeofencesOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm geofenceRepositoryMiddleware) Update(ctx context.Context, th lbs.GeofenceRecord) error {
	span := createSpan(ctx, trm.tracer, updateGeofenceOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm geofenceRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (lbs.GeofenceRecord, error) {
	span := createSpan(ctx, trm.tracer, retrieveGeofenceByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm geofenceRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata lbs.Metadata) (lbs.GeofencesPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllGeofencesOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm geofenceRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeGeofenceOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type geofenceCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  lbs.GeofenceCache
}

// GeofenceCacheMiddleware tracks request and their latency, and adds spans
// to context.
func GeofenceCacheMiddleware(tracer opentracing.Tracer, cache lbs.GeofenceCache) lbs.GeofenceCache {
	return geofenceCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm geofenceCacheMiddleware) Save(ctx context.Context, geofenceKey string, geofenceID string) error {
	span := createSpan(ctx, tcm.tracer, saveGeofenceOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, geofenceKey, geofenceID)
}

func (tcm geofenceCacheMiddleware) ID(ctx context.Context, geofenceKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveGeofenceIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, geofenceKey)
}

func (tcm geofenceCacheMiddleware) Remove(ctx context.Context, geofenceID string) error {
	span := createSpan(ctx, tcm.tracer, removeGeofenceOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, geofenceID)
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
