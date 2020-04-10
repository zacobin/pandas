// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package tracing contains middlewares that will add spans
// to existing traces.
package tracing

import (
	"context"

	"github.com/cloustone/pandas/authz"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveOp             = "save_op"
	retrieveByIDOp     = "retrieve_by_id"
	generateResetToken = "generate_reset_token"
	updatePassword     = "update_password"
	sendPasswordReset  = "send_reset_password"
	revokeRole         = "revoke_role"
	listRole           = "list_role"
)

var _ authz.RoleRepository = (*roleRepositoryMiddleware)(nil)

type roleRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   authz.RoleRepository
}

// RoleRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func RoleRepositoryMiddleware(repo authz.RoleRepository, tracer opentracing.Tracer) authz.RoleRepository {
	return roleRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (rrm roleRepositoryMiddleware) Save(ctx context.Context, role authz.Role) error {
	span := createSpan(ctx, rrm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Save(ctx, role)
}

func (rrm roleRepositoryMiddleware) Update(ctx context.Context, role authz.Role) error {
	span := createSpan(ctx, rrm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Update(ctx, role)
}

func (rrm roleRepositoryMiddleware) Retrieve(ctx context.Context, name string) (authz.Role, error) {
	span := createSpan(ctx, rrm.tracer, retrieveByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return rrm.repo.Retrieve(ctx, name)
}

func (urm roleRepositoryMiddleware) Revoke(ctx context.Context, name string) error {
	span := createSpan(ctx, urm.tracer, revokeRole)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return urm.repo.Revoke(ctx, name)
}

func (urm roleRepositoryMiddleware) List(ctx context.Context) ([]authz.Role, error) {
	span := createSpan(ctx, urm.tracer, listRole)
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
