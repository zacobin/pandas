// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/cloustone/pandas/mainflux"
	log "github.com/cloustone/pandas/pkg/logger"
	"github.com/cloustone/pandas/rulechain"
)

var _ rulechain.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    rulechain.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc rulechain.Service, logger log.Logger) rulechain.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) AddNewRuleChain(ctx context.Context, rulechain rulechain.RuleChain) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method AddNewRuleChain for rulechain %s took %s to complete", rulechain.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AddNewRuleChain(ctx, rulechain)
}

func (lm *loggingMiddleware) GetRuleChainInfo(ctx context.Context, UserID string, RuleChainID string) (rulechain rulechain.RuleChain, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method getrulechaininfo for rulechain %s took %s to complete", RuleChainID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.GetRuleChainInfo(ctx, UserID, RuleChainID)
}

func (lm *loggingMiddleware) UpdateRuleChain(ctx context.Context, rulechain rulechain.RuleChain) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method updaterulechain for rulechain %s took %s to complete", rulechain.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateRuleChain(ctx, rulechain)
}

func (lm *loggingMiddleware) RevokeRuleChain(ctx context.Context, UserID string, RuleChainID string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method revokerulechain for rulechain %s took %s to complete", RuleChainID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RevokeRuleChain(ctx, UserID, RuleChainID)
}

func (lm *loggingMiddleware) ListRuleChain(ctx context.Context, UserID string) (rulechains []rulechain.RuleChain, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method listrulechain  took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListRuleChain(ctx, UserID)
}

func (lm *loggingMiddleware) StartRuleChain(ctx context.Context, UserID string, RuleChainID string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method startrulechain for rulechain %s took %s to complete", RuleChainID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.StartRuleChain(ctx, UserID, RuleChainID)
}

func (lm *loggingMiddleware) StopRuleChain(ctx context.Context, UserID string, RuleChainID string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method stoprulechain for rulechain %s took %s to complete", RuleChainID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.StopRuleChain(ctx, UserID, RuleChainID)
}

func (lm *loggingMiddleware) SaveStates(msg *mainflux.Message) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method savesates took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.SaveStates(msg)
}
