// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/cloustone/pandas/pms"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveProjectOp               = "save_project"
	saveProjectsOp              = "save_project"
	updateProjectOp             = "update_project"
	updateProjectKeyOp          = "update_project_by_key"
	retrieveProjectByIDOp       = "retrieve_project_by_id"
	retrieveProjectByKeyOp      = "retrieve_project_by_key"
	retrieveAllProjectsOp       = "retrieve_all_projects"
	retrieveProjectsByChannelOp = "retrieve_pms_by_chan"
	removeProjectOp             = "remove_project"
	retrieveProjectIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ pms.ProjectRepository = (*projectRepositoryMiddleware)(nil)
	_ pms.ProjectCache      = (*projectCacheMiddleware)(nil)
)

type projectRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   pms.ProjectRepository
}

// ProjectRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func ProjectRepositoryMiddleware(tracer opentracing.Tracer, repo pms.ProjectRepository) pms.ProjectRepository {
	return projectRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm projectRepositoryMiddleware) Save(ctx context.Context, ths ...pms.Project) ([]pms.Project, error) {
	span := createSpan(ctx, trm.tracer, saveProjectsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm projectRepositoryMiddleware) Update(ctx context.Context, th pms.Project) error {
	span := createSpan(ctx, trm.tracer, updateProjectOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm projectRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (pms.Project, error) {
	span := createSpan(ctx, trm.tracer, retrieveProjectByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm projectRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata pms.Metadata) (pms.ProjectsPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllProjectsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm projectRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeProjectOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type projectCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  pms.ProjectCache
}

// ProjectCacheMiddleware tracks request and their latency, and adds spans
// to context.
func ProjectCacheMiddleware(tracer opentracing.Tracer, cache pms.ProjectCache) pms.ProjectCache {
	return projectCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm projectCacheMiddleware) Save(ctx context.Context, projectKey string, projectID string) error {
	span := createSpan(ctx, tcm.tracer, saveProjectOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, projectKey, projectID)
}

func (tcm projectCacheMiddleware) ID(ctx context.Context, projectKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveProjectIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, projectKey)
}

func (tcm projectCacheMiddleware) Remove(ctx context.Context, projectID string) error {
	span := createSpan(ctx, tcm.tracer, removeProjectOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, projectID)
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
