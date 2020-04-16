// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"time"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/rulechain"
	"github.com/go-kit/kit/metrics"
)

var _ rulechain.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     rulechain.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc rulechain.Service, counter metrics.Counter, latency metrics.Histogram) rulechain.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) AddNewRuleChain(ctx context.Context, token string, rulechain rulechain.RuleChain) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "addnewrulechain").Add(1)
		ms.latency.With("method", "addnewrulechain").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddNewRuleChain(ctx, token, rulechain)
}

func (ms *metricsMiddleware) GetRuleChainInfo(ctx context.Context, token string, RuleChainID string) (rulechain.RuleChain, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "getrulechaininfo").Add(1)
		ms.latency.With("method", "getrulechaininfo").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetRuleChainInfo(ctx, token, RuleChainID)
}

func (ms *metricsMiddleware) UpdateRuleChain(ctx context.Context, token string, rulechain rulechain.RuleChain) (rulechain.RuleChain, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "updaterulechain").Add(1)
		ms.latency.With("method", "updaterulechain").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateRuleChain(ctx, token, rulechain)
}

func (ms *metricsMiddleware) RevokeRuleChain(ctx context.Context, token string, RuleChainID string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "revokerulechain").Add(1)
		ms.latency.With("method", "revokerulechain").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RevokeRuleChain(ctx, token, RuleChainID)
}

func (ms *metricsMiddleware) ListRuleChain(ctx context.Context, token string, offset uint64, limit uint64) (rulechain.RuleChainPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "listrulechain").Add(1)
		ms.latency.With("method", "listrulechain").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListRuleChain(ctx, token, offset, limit)
}

func (ms *metricsMiddleware) UpdateRuleChainStatus(ctx context.Context, token string, RuleChainID string, updatestatus string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "startrulechain").Add(1)
		ms.latency.With("method", "startrulechain").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateRuleChainStatus(ctx, token, RuleChainID, updatestatus)
}

func (ms *metricsMiddleware) SaveStates(msg *mainflux.Message) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "savestates").Add(1)
		ms.latency.With("method", "savestates").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.SaveStates(msg)
}
