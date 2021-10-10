package kuiper

import "context"

// Rule represents a Mainflux thing. Each thing is owned by one user, and
// it is assigned with the unique identifier and (temporary) access key.
type Rule struct {
	Owner    string
	ID       string
	Name     string
	SQL      string
	Metadata Metadata
}

// RulesPage contains page related metadata as well as list of things that
// belong to this page.
type RulesPage struct {
	PageMetadata
	Rules []Rule
}

// RuleRepository specifies a thing persistence API.
type RuleRepository interface {
	// Save persists multiple things. Rules are saved using a transaction. If one thing
	// fails then none will be saved. Successful operation is indicated by non-nil
	// error response.
	Save(context.Context, ...Rule) ([]Rule, error)

	// Update performs an update to the existing thing. A non-nil error is
	// returned to indicate operation failure.
	Update(context.Context, Rule) error

	// RetrieveByID retrieves the rule having the provided identifier, that is owned
	// by the specified user.
	RetrieveByID(context.Context, string, string) (Rule, error)

	// RetrieveAll retrieves the subset of things owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (RulesPage, error)

	// Remove removes the thing having the provided identifier, that is owned
	// by the specified user.
	Remove(context.Context, string, string) error
}

// RuleCache contains thing caching interface.
type RuleCache interface {
	// Save stores pair thing key, thing id.
	Save(context.Context, string, string) error

	// ID returns thing ID for given key.
	ID(context.Context, string) (string, error)

	// Removes thing from cache.
	Remove(context.Context, string) error
}
