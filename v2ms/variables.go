package v2ms

import (
	"context"
	"time"
)

// Metadata stores arbitrary variable data
type Metadata map[string]interface{}

// Variable represents a variable. Each variable is binded iwth one thing, and
// is assigned with the unique identifier.
type Variable struct {
	Owner          string
	ID             string
	Name           string
	ThingID        string
	ThingAttribute string
	Created        time.Time
	Updated        time.Time
	Revision       int
	Metadata       Metadata
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
	Name   string
}

// VariablesPage contains page related metadata as well as a list of variables that
// belong to this page.
type VariablesPage struct {
	PageMetadata
	Variables []Variable
}

// VariableRepository specifies a variable persistence API.
type VariableRepository interface {
	// Save persists the variable
	Save(context.Context, ...Variable) ([]Variable, error)

	// Update performs an update to the existing variable. A non-nil error is
	// returned to indicate operation failure.
	Update(context.Context, Variable) error

	// RetrieveByID retrieves the variable having the provided identifier.
	RetrieveByID(ctx context.Context, owner, id string) (Variable, error)

	// RetrieveAll retrieves the subset of things owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (VariablesPage, error)

	// Remove removes the variable having the provided identifier.
	Remove(ctx context.Context, owner, id string) error
}

// VariableCache contains thing caching interface.
type VariableCache interface {
	// Save stores pair thing key, thing id.
	Save(context.Context, string, string) error

	// ID returns thing ID for given key.
	ID(context.Context, string) (string, error)

	// Removes thing from cache.
	Remove(context.Context, string) error
}
