// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package tracing contains middlewares that will add spans
// to existing traces.
package tracing

import (
	"context"

	"github.com/cloustone/pandas/realms"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveOp             = "save_op"
	retrieveByIDOp     = "retrieve_by_id"
	generateResetToken = "generate_reset_token"
	updatePassword     = "update_password"
	sendPasswordReset  = "send_reset_password"
	revokeRealm        = "revoke_realm"
	listRealm          = "list_realm"
)

var _ realms.RealmRepository = (*realmRepositoryMiddleware)(nil)

type realmRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   realms.RealmRepository
}

// RealmRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func RealmRepositoryMiddleware(repo realms.RealmRepository, tracer opentracing.Tracer) realms.RealmRepository {
	return realmRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (urm realmRepositoryMiddleware) Save(ctx context.Context, realm realms.Realm) error {
	span := createSpan(ctx, urm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Save(ctx, realm)
}

func (urm realmRepositoryMiddleware) Update(ctx context.Context, realm realms.Realm) error {
	span := createSpan(ctx, urm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Update(ctx, realm)
}

func (urm realmRepositoryMiddleware) Retrieve(ctx context.Context, name string) (realms.Realm, error) {
	span := createSpan(ctx, urm.tracer, retrieveByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Retrieve(ctx, name)
}

func (urm realmRepositoryMiddleware) Revoke(ctx context.Context, name string) error {
	span := createSpan(ctx, urm.tracer, revokeRealm)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Revoke(ctx, name)
}

func (urm realmRepositoryMiddleware) List(ctx context.Context) ([]realms.Realm, error) {
	span := createSpan(ctx, urm.tracer, listRealm)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.List(ctx)
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
