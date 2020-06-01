// SPDX-License-Identifier: Apache-2.0

package alerts

import (
	"context"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/errors"
)

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict = errors.New("email already taken")

	// ErrMalformedEntity indicates malformed entity specification
	// (e.g. invalid alertname or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrAlertNotFound indicates a non-existent alert request.
	ErrAlertNotFound = errors.New("non-existent alert")

	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = errors.New("Failed to scan metadata")

	// ErrMissingEmail indicates missing email for password reset request.
	ErrMissingEmail = errors.New("missing email for password reset")

	// ErrUnauthorizedPrincipal indicate the pricipal can not be recognized
	ErrUnauthorizedPrincipal = errors.New("unauthorized principal")
)

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
	Name   string
}

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateAlert creates new alert . In case of the failed registration, a
	// non-nil error value is returned.
	CreateAlert(ctx context.Context, token string, alert Alert) (Alert, error)

	// RetrieveAlert return authenticated alert info for the given token.
	RetrieveAlert(ctx context.Context, token, name string) (Alert, error)

	// UpdateAlert updates the alert metadata.
	UpdateAlert(ctx context.Context, token string, alert Alert) error

	// RevokeAlert remove a alert
	RevokeAlert(ctx context.Context, token, name string) error

	// RetrieveAlerts retrieves the subset of alerts owned by the specified user.
	RetrieveAlerts(context.Context, string, uint64, uint64, string, Metadata) (AlertsPage, error)

	// CreateAlertRule creates new alert rule. In case of the failed registration, a
	// non-nil error value is returned.
	CreateAlertRule(ctx context.Context, token string, alert AlertRule) (AlertRule, error)

	// RetrieveAlertRule authenticated alert rule info for the given token.
	RetrieveAlertRule(ctx context.Context, token, name string) (AlertRule, error)

	// UpdateAlert updates the alert metadata.
	UpdateAlertRule(ctx context.Context, token string, alert AlertRule) error

	// RevokeAlertRule remove a alert
	RevokeAlertRule(ctx context.Context, token, name string) error

	// RetrieveAlerts retrieves the subset of alerts owned by the specified user.
	RetrieveAlertRules(context.Context, string, uint64, uint64, string, Metadata) (AlertRulesPage, error)
}

var _ Service = (*alertService)(nil)

type alertService struct {
	hasher Hasher
	auth   mainflux.AuthNServiceClient
	alerts AlertRepository
	alarms AlarmRepository
	rules  AlertRuleRepository
	idp    IdentityProvider
}

// New instantiates the alerts service implementation
func New(auth mainflux.AuthNServiceClient, hasher Hasher, idp IdentityProvider, alerts AlertRepository, alarms AlarmRepository, rules AlertRuleRepository) Service {
	return &alertService{
		hasher: hasher,
		auth:   auth,
		alerts: alerts,
		alarms: alarms,
		rules:  rules,
		idp:    idp,
	}
}

func (svc alertService) CreateAlert(ctx context.Context, token string, alert Alert) (Alert, error) {
	owner, err := svc.identify(ctx, token)
	if err != nil {
		return Alert{}, err
	}
	id, err := svc.idp.ID()
	if err != nil {
		return Alert{}, err
	}
	alert.ID = id
	alert.Owner = owner
	return svc.alerts.Save(ctx, alert)
}

func (svc alertService) RetrieveAlert(ctx context.Context, token, name string) (Alert, error) {
	owner, err := svc.identify(ctx, token)
	if err != nil {
		return Alert{}, err
	}

	alert, err := svc.alerts.Retrieve(ctx, owner, name)
	if err != nil {
		return Alert{}, errors.Wrap(ErrUnauthorizedAccess, err)
	}

	return alert, nil
}

func (svc alertService) UpdateAlert(ctx context.Context, token string, alert Alert) error {
	user, err := svc.identify(ctx, token)
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	if alert.Owner != user {
		return ErrMalformedEntity
	}
	return svc.alerts.Update(ctx, alert)
}

func (svc alertService) RetrieveAlerts(ctx context.Context, token string, offset uint64, limit uint64, id string, meta Metadata) (AlertsPage, error) {
	user, err := svc.identify(ctx, token)
	if err != nil {
		return AlertsPage{}, errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.alerts.RetrieveAll(ctx, user, offset, limit, id, meta)
}

func (svc alertService) RevokeAlert(ctx context.Context, token string, name string) error {
	user, err := svc.identify(ctx, token)
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.alerts.Revoke(ctx, user, name)
}

func (svc alertService) CreateAlertRule(ctx context.Context, token string, alertRule AlertRule) (AlertRule, error) {
	user, err := svc.identify(ctx, token)
	if err != nil {
		return AlertRule{}, err
	}
	id, err := svc.idp.ID()
	if err != nil {
		return AlertRule{}, err
	}
	alertRule.ID = id
	alertRule.Owner = user
	return svc.rules.Save(ctx, alertRule)
}

func (svc alertService) RetrieveAlertRule(ctx context.Context, token, name string) (AlertRule, error) {
	user, err := svc.identify(ctx, token)
	if err != nil {
		return AlertRule{}, err
	}

	alertRule, err := svc.rules.Retrieve(ctx, user, name)
	if err != nil {
		return AlertRule{}, errors.Wrap(ErrUnauthorizedAccess, err)
	}

	return alertRule, nil
}

func (svc alertService) UpdateAlertRule(ctx context.Context, token string, alertRule AlertRule) error {
	user, err := svc.identify(ctx, token)
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	if alertRule.Owner != user {
		return ErrMalformedEntity
	}
	return svc.rules.Update(ctx, alertRule)
}

func (svc alertService) RetrieveAlertRules(ctx context.Context, token string, offset uint64, limit uint64, rule string, meta Metadata) (AlertRulesPage, error) {
	user, err := svc.identify(ctx, token)
	if err != nil {
		return AlertRulesPage{}, errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.rules.RetrieveAll(ctx, user, offset, limit, rule, meta)
}

func (svc alertService) RevokeAlertRule(ctx context.Context, token string, name string) error {
	user, err := svc.identify(ctx, token)
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.rules.Revoke(ctx, user, name)
}

func (svc alertService) identify(ctx context.Context, token string) (string, error) {
	owner, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return owner.GetValue(), nil
}
