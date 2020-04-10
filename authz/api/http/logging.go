// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package httpapi

import (
	"context"
	"fmt"
	"time"

	"github.com/cloustone/pandas/authz"
	log "github.com/cloustone/pandas/pkg/logger"
)

var _ authz.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    authz.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc authz.Service, logger log.Logger) authz.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) RetrieveRole(ctx context.Context, token, roleName string) (role authz.Role, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method retrieve for role %s took %s to complete", roleName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.RetrieveRole(ctx, token, roleName)
}

func (lm *loggingMiddleware) ListRoles(ctx context.Context, token string) (roles []authz.Role, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_roles took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListRoles(ctx, token)
}

func (lm *loggingMiddleware) UpdateRole(ctx context.Context, token string, role authz.Role) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_role for role %s took %s to complete", role.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateRole(ctx, token, role)
}

func (lm *loggingMiddleware) Authorize(ctx context.Context, token string, roleName string, subject authz.Subject) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method authorize for role %s took %s to complete", roleName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Authorize(ctx, token, roleName, subject)
}
