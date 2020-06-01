// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/cloustone/pandas/alerts"
	log "github.com/cloustone/pandas/pkg/logger"
)

var _ alerts.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    alerts.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc alerts.Service, logger log.Logger) alerts.Service {
	return &loggingMiddleware{logger, svc}
}

// Alert
func (lm *loggingMiddleware) CreateAlert(ctx context.Context, token string, alert alerts.Alert) (newAlert alerts.Alert, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for alert %s took %s to complete", alert.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.CreateAlert(ctx, token, alert)
}

func (lm *loggingMiddleware) RetrieveAlert(ctx context.Context, token, alertName string) (alert alerts.Alert, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method alert_info %s took %s to complete", alertName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RetrieveAlert(ctx, token, alertName)
}

func (lm *loggingMiddleware) UpdateAlert(ctx context.Context, token string, alert alerts.Alert) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_alert for alert %s took %s to complete", alert.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateAlert(ctx, token, alert)
}

func (lm *loggingMiddleware) RevokeAlert(ctx context.Context, token string, alertName string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method revoke_alert for alert %s took %s to complete", alertName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RevokeAlert(ctx, token, alertName)
}

func (lm *loggingMiddleware) RetrieveAlerts(ctx context.Context, token string, offset, limit uint64, name string, metadata alerts.Metadata) (page alerts.AlertsPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_alerts  took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RetrieveAlerts(ctx, token, offset, limit, name, metadata)
}

/*
func (lm loggingMiddleware) Revoke(ctx context.Context, token, name string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method revoke_alerts  took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Revoke(ctx, token, name)

}
*/

// AlertRule
func (lm *loggingMiddleware) CreateAlertRule(ctx context.Context, token string, alertRule alerts.AlertRule) (newAlertRule alerts.AlertRule, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for alert rule %s took %s to complete", alertRule.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.CreateAlertRule(ctx, token, alertRule)
}

func (lm *loggingMiddleware) RetrieveAlertRule(ctx context.Context, token string, alertRuleName string) (alert alerts.AlertRule, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method alert_rule_info %s took %s to complete", alertRuleName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RetrieveAlertRule(ctx, token, alertRuleName)
}

func (lm *loggingMiddleware) UpdateAlertRule(ctx context.Context, token string, alertRule alerts.AlertRule) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_alert for alert rule %s took %s to complete", alertRule.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateAlertRule(ctx, token, alertRule)
}

func (lm *loggingMiddleware) RevokeAlertRule(ctx context.Context, token string, alertRuleName string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method revoke_alert for alert rule %s took %s to complete", alertRuleName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RevokeAlert(ctx, token, alertRuleName)
}

func (lm *loggingMiddleware) RetrieveAlertRules(ctx context.Context, token string, offset, limit uint64, name string, meta alerts.Metadata) (page alerts.AlertRulesPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_alerts took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RetrieveAlertRules(ctx, token, offset, limit, name, meta)
}
