// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/cloustone/pandas/pkg/logger"
	swagger "github.com/cloustone/pandas/swagger"
)

var _ swagger.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    swagger.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc swagger.Service, logger log.Logger) swagger.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) RetrieveSwaggerConfigs(ctx context.Context, token string) (viewed swagger.Configs, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_configs for token %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RetrieveSwaggerConfigs(ctx, token)
}

func (lm *loggingMiddleware) RetrieveDownstreamSwagger(ctx context.Context, token, module string) (viewed swagger.DownstreamSwagger, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_swagger for token %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RetrieveDownstreamSwagger(ctx, token, module)
}
