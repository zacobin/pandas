// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/pms"
)

var _ pms.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    pms.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc pms.Service, logger log.Logger) pms.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) AddProject(ctx context.Context, token string, project pms.Project) (saved pms.Project, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_project for token %s and project %s took %s to complete", token, project.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.AddProject(ctx, token, project)
}

func (lm *loggingMiddleware) UpdateProject(ctx context.Context, token string, project pms.Project) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_project for token %s and view %s took %s to complete", token, project.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateProject(ctx, token, project)
}

func (lm *loggingMiddleware) ViewProject(ctx context.Context, token, id string) (viewed pms.Project, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_project for token %s and view %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewProject(ctx, token, id)
}

func (lm *loggingMiddleware) ListProjects(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata pms.Metadata) (tw pms.ProjectsPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_projects for token %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListProjects(ctx, token, offset, limit, name, metadata)
}

func (lm *loggingMiddleware) RemoveProject(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_project for token %s and view %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveProject(ctx, token, id)
}
