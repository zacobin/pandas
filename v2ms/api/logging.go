// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/v2ms"
)

var _ v2ms.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    v2ms.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc v2ms.Service, logger log.Logger) v2ms.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) AddView(ctx context.Context, token string, view v2ms.View) (saved v2ms.View, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_view for token %s and view %s took %s to complete", token, view.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.AddView(ctx, token, view)
}

func (lm *loggingMiddleware) UpdateView(ctx context.Context, token string, view v2ms.View) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_view for token %s and view %s took %s to complete", token, view.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateView(ctx, token, view)
}

func (lm *loggingMiddleware) ViewView(ctx context.Context, token, id string) (viewed v2ms.View, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_view for token %s and view %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewView(ctx, token, id)
}

func (lm *loggingMiddleware) ListViews(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata v2ms.Metadata) (tw v2ms.ViewsPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_views for token %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListViews(ctx, token, offset, limit, name, metadata)
}

func (lm *loggingMiddleware) RemoveView(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_view for token %s and view %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveView(ctx, token, id)
}

// Variable
func (lm *loggingMiddleware) AddVariable(ctx context.Context, token string, variable v2ms.Variable) (saved v2ms.Variable, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_variable for token %s and variable %s took %s to complete", token, variable.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.AddVariable(ctx, token, variable)
}

func (lm *loggingMiddleware) UpdateVariable(ctx context.Context, token string, variable v2ms.Variable) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_variable for token %s and view %s took %s to complete", token, variable.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateVariable(ctx, token, variable)
}

func (lm *loggingMiddleware) ViewVariable(ctx context.Context, token, id string) (viewed v2ms.Variable, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_variable for token %s and view %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewVariable(ctx, token, id)
}

func (lm *loggingMiddleware) ListVariables(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata v2ms.Metadata) (tw v2ms.VariablesPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_variables for token %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListVariables(ctx, token, offset, limit, name, metadata)
}

func (lm *loggingMiddleware) RemoveVariable(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_variable for token %s and view %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveVariable(ctx, token, id)
}
