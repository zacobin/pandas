// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package v2ms

import (
	"context"
	"time"
)

//  View represents a view thant conatains a set of variables and diplayable
//  elements. Each view is assigned with the unique identifier.
type View struct {
	Owner    string
	ID       string
	Name     string
	Created  time.Time
	Updated  time.Time
	Revision int
	Metadata Metadata
}

// ViewsPage contains page related metadata as well as a list of views that
// belong to this page.
type ViewsPage struct {
	PageMetadata
	Views []View
}

// ViewRepository specifies a variable persistence API.
type ViewRepository interface {
	// Save persists the view
	Save(context.Context, ...View) ([]View, error)

	// Update performs an update to the existing view. A non-nil error is
	// returned to indicate operation failure.
	Update(context.Context, View) error

	// RetrieveByID retrieves the view having the provided identifier.
	RetrieveByID(ctx context.Context, owner, id string) (View, error)

	// RetrieveAll retrieves the subset of things owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (ViewsPage, error)

	// Remove removes the view having the provided identifier.
	Remove(ctx context.Context, owner, id string) error
}

// ViewCache contains thing caching interface.
type ViewCache interface {
	// Save stores pair thing key, thing id.
	Save(context.Context, string, string) error

	// ID returns thing ID for given key.
	ID(context.Context, string) (string, error)

	// Removes thing from cache.
	Remove(context.Context, string) error
}
