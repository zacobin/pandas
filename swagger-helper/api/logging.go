// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/cloustone/pandas/pkg/logger"
	swagger_helper "github.com/cloustone/pandas/swagger-helper"
)

var _ swagger_helper.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    swagger_helper.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc swagger_helper.Service, logger log.Logger) swagger_helper.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) RetrieveDownstreamSwagger(ctx context.Context, token, module string) (viewed swagger_helper.DownstreamSwagger, err error) {
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
