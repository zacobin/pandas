package pms

import (
	"context"
	"time"
)

// Metadata stores arbitrary variable data
type Metadata map[string]interface{}

// Project represents a variable. Each variable is binded iwth one thing, and
// is assigned with the unique identifier.
type Project struct {
	Owner    string
	ID       string
	Name     string
	Created  time.Time
	Updated  time.Time
	Revision int
	Metadata Metadata
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
	Name   string
}

// ProjectsPage contains page related metadata as well as a list of variables that
// belong to this page.
type ProjectsPage struct {
	PageMetadata
	Projects []Project
}

// ProjectRepository specifies a variable persistence API.
type ProjectRepository interface {
	// Save persists the variable
	Save(context.Context, ...Project) ([]Project, error)

	// Update performs an update to the existing variable. A non-nil error is
	// returned to indicate operation failure.
	Update(context.Context, Project) error

	// RetrieveByID retrieves the variable having the provided identifier.
	RetrieveByID(ctx context.Context, owner, id string) (Project, error)

	// RetrieveAll retrieves the subset of things owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (ProjectsPage, error)

	// Remove removes the variable having the provided identifier.
	Remove(ctx context.Context, owner, id string) error
}

// ProjectCache contains thing caching interface.
type ProjectCache interface {
	// Save stores pair thing key, thing id.
	Save(context.Context, string, string) error

	// ID returns thing ID for given key.
	ID(context.Context, string) (string, error)

	// Removes thing from cache.
	Remove(context.Context, string) error
}
