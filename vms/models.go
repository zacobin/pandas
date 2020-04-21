package vms

import (
	"context"
	"time"
)

// Model represents a model. Each model is binded iwth one thing, and
// is assigned with the unique identifier.
type Model struct {
	Owner    string
	ID       string
	Name     string
	Created  time.Time
	Updated  time.Time
	Revision int
	Metadata Metadata
}

// ModelsPage contains page related metadata as well as a list of models that
// belong to this page.
type ModelsPage struct {
	PageMetadata
	Models []Model
}

// ModelRepository specifies a model persistence API.
type ModelRepository interface {
	// Save persists the model
	Save(context.Context, ...Model) ([]Model, error)

	// Update performs an update to the existing model. A non-nil error is
	// returned to indicate operation failure.
	Update(context.Context, Model) error

	// Retrieve retrieves the model having the provied identifier without
	// owner, using only internal
	Retrieve(ctx context.Context, id string) (Model, error)

	// RetrieveByID retrieves the model having the provided identifier.
	RetrieveByID(ctx context.Context, owner, id string) (Model, error)

	// RetrieveAll retrieves the subset of things owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (ModelsPage, error)

	// Remove removes the model having the provided identifier.
	Remove(ctx context.Context, owner, id string) error
}

// ModelCache contains thing caching interface.
type ModelCache interface {
	// Save stores pair thing key, thing id.
	Save(context.Context, string, string) error

	// ID returns thing ID for given key.
	ID(context.Context, string) (string, error)

	// Removes thing from cache.
	Remove(context.Context, string) error
}
