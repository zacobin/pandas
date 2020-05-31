// SPDX-License-Identifier: Apache-2.0

// Package tracing contains middlewares that will add spans
// to existing traces.
package tracing

import (
	"context"

	"github.com/cloustone/pandas/alarms"
	"github.com/cloustone/pandas/alerts"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveOp     = "save_op"
	retrieveOp = "retrieve_op"
	updateOp   = "update_op"
	revokeOp   = "revoke_op"
	listOp     = "list_op"
)

var _ alerts.AlarmRepository = (*alarmRepositoryMiddleware)(nil)

type alarmRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   alarms.AlarmRepository
}

// AlarmRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func AlarmRepositoryMiddleware(repo alarms.AlarmRepository, tracer opentracing.Tracer) alarms.AlarmRepository {
	return alarmRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (arm alarmRepositoryMiddleware) Save(ctx context.Context, alarm ...alarms.Alarm) ([]alarms.Alarm, error) {
	span := createSpan(ctx, arm.tracer, saveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.Save(ctx, alarm...)
}

func (arm alarmRepositoryMiddleware) Update(ctx context.Context, alarm alarms.Alarm) error {
	span := createSpan(ctx, arm.tracer, updageOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.Update(ctx, alarm)
}

func (arm alarmRepositoryMiddleware) Retrieve(ctx context.Context, owner, name string) (alarms.Alarm, error) {
	span := createSpan(ctx, arm.tracer, retrieveOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.Retrieve(ctx, owner, name)
}

func (arm alarmRepositoryMiddleware) Revoke(ctx context.Context, owner, name string) error {
	span := createSpan(ctx, arm.tracer, revokeOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.Revoke(ctx, owner, name)
}

func (arm alarmRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, meta Metadata) (alerts.AlarmsPage, error) {
	span := createSpan(ctx, arm.tracer, listOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return arm.repo.RetrieveAll(ctx, owner, offset, limit, name, meta)
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
