// SPDX-License-Identifier: Apache-2.0
package alerts

import (
	"context"
)

// AlertsPage contains page related metadata as well as a list of variables that
// belong to this page.
type AlertsPage struct {
	PageMetadata
	Alerts []Alert
}

// AlertRepository specifies alert persistence API
type AlertRepository interface {
	// Save persists the alert
	Save(context.Context, ...Alert) ([]Alert, error)

	// Update updates the alert metdata
	Update(context.Context, Alert) error

	// Retrieve return alert by its identifier (i.e name)
	Retrieve(context.Context, string, string) (Alert, error)

	// RevokeAlert remove alert
	Revoke(context.Context, string, string) error

	// RetrieveAll retrieves the subset of alerts owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (AlertsPage, error)
}
