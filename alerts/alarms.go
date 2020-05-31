// SPDX-License-Identifier: Apache-2.0
package alerts

import "context"

// AlarmsPage contains page related metadata as well as a list of variables that
// belong to this page.
type AlarmsPage struct {
	PageMetadata
	Alarms []Alarm
}

// AlarmRepository specifies alert persistence API
type AlarmRepository interface {
	// Save persists the alert
	Save(context.Context, ...Alarm) ([]Alarm, error)

	// Update updates the alert metdata
	Update(context.Context, Alarm) error

	// Retrieve return alert by its identifier (i.e name)
	Retrieve(context.Context, string, string) (Alarm, error)

	// RevokeAlarm remove alert
	Revoke(context.Context, string, string) error

	// RetrieveAll retrieves the subset of alarms owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (AlarmsPage, error)
}
