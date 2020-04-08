// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/realms"
)

var _ realms.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    realms.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc realms.Service, logger log.Logger) realms.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) Register(ctx context.Context, realm realms.Realm) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for realm %s took %s to complete", realm.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Register(ctx, realm)
}

func (lm *loggingMiddleware) RealmInfo(ctx context.Context, token, realmName string) (realm realms.Realm, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method realmi_info %s took %s to complete", realmName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RealmInfo(ctx, token, realmName)
}

func (lm *loggingMiddleware) UpdateRealm(ctx context.Context, token string, realm realms.Realm) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_realm for realm %s took %s to complete", realm.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateRealm(ctx, token, realm)
}

func (lm *loggingMiddleware) RevokeRealm(ctx context.Context, token string, realmName string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method revoke_realm for realm %s took %s to complete", realmName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RevokeRealm(ctx, token, realmName)
}

func (lm *loggingMiddleware) ListRealms(ctx context.Context, token string) (realms []realms.Realm, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_realms  took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListRealms(ctx, token)
}

func (lm *loggingMiddleware) Identify(ctx context.Context, token string, principal realms.Principal) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method identify_principal for principal %s took %s to complete", principal.Username, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Identify(ctx, token, principal)
}
