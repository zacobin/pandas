// SPDX-License-Identifier: Apache-2.0
package alerts

import (
	"context"
)

// AlertRulesPage contains page related metadata as well as a list of variables that
// belong to this page.
type AlertRulesPage struct {
	PageMetadata
	AlertRules []AlertRule
}

// AlertRuleRepository specifies alert persistence API
type AlertRuleRepository interface {
	// Save persists the alert rule
	Save(context.Context, AlertRule) (AlertRule, error)

	// Update updates the alert metdata
	Update(context.Context, AlertRule) error

	// Retrieve return alert by its identifier (i.e name)
	Retrieve(context.Context, string, string) (AlertRule, error)

	// Revoke remove alert rule
	Revoke(context.Context, string, string) error

	// RetrieveAll retrieves the subset of alert rule owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (AlertRulesPage, error)
}
