// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"time"

	"github.com/cloustone/pandas/alerts"
	"github.com/go-kit/kit/metrics"
)

var _ alerts.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     alerts.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc alerts.Service, counter metrics.Counter, latency metrics.Histogram) alerts.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

// Alert
func (ms *metricsMiddleware) CreateAlert(ctx context.Context, token string, alert alerts.Alert) (alerts.Alert, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_alert").Add(1)
		ms.latency.With("method", "create_alert").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateAlert(ctx, token, alert)
}

func (ms *metricsMiddleware) RetrieveAlert(ctx context.Context, token, name string) (alerts.Alert, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve_alert").Add(1)
		ms.latency.With("method", "retrieve_alert").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveAlert(ctx, token, name)
}

func (ms *metricsMiddleware) UpdateAlert(ctx context.Context, token string, alert alerts.Alert) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_alert").Add(1)
		ms.latency.With("method", "update_alert").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateAlert(ctx, token, alert)
}

func (ms *metricsMiddleware) RevokeAlert(ctx context.Context, token string, name string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "revoke_alert").Add(1)
		ms.latency.With("method", "revoke_alert").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RevokeAlert(ctx, token, name)
}

func (ms *metricsMiddleware) RetrieveAlerts(ctx context.Context, token string, offset, limit uint64, name string, metadata alerts.Metadata) (alerts.AlertsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_alerts").Add(1)
		ms.latency.With("method", "list_alerts").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveAlerts(ctx, token, offset, limit, name, metadata)
}

// Alarms
func (ms *metricsMiddleware) CreateAlertRule(ctx context.Context, token string, alertRule alerts.AlertRule) (alerts.AlertRule, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_alert_rule").Add(1)
		ms.latency.With("method", "create_alert_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateAlertRule(ctx, token, alertRule)
}

func (ms *metricsMiddleware) RetrieveAlertRule(ctx context.Context, token, name string) (alerts.AlertRule, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve_alert_rule").Add(1)
		ms.latency.With("method", "retrieve_alert_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveAlertRule(ctx, token, name)
}

func (ms *metricsMiddleware) UpdateAlertRule(ctx context.Context, token string, alertRule alerts.AlertRule) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_alert_rule").Add(1)
		ms.latency.With("method", "update_alert_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateAlertRule(ctx, token, alertRule)
}

func (ms *metricsMiddleware) RevokeAlertRule(ctx context.Context, token string, name string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "revoke_alert_rule").Add(1)
		ms.latency.With("method", "revoke_alert_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RevokeAlertRule(ctx, token, name)
}

func (ms *metricsMiddleware) RetrieveAlertRules(ctx context.Context, token string, offset, limit uint64, name string, meta alerts.Metadata) (alerts.AlertRulesPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_alert_rules").Add(1)
		ms.latency.With("method", "list_alert_rules").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveAlertRules(ctx, token, offset, limit, name, meta)
}
