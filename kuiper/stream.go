package kuiper

import "context"

// Metadata to be used for mainflux stream or stream for customized
// describing of particular stream
type Metadata map[string]interface{}

// Stream represents a stream endpoint. Each stream is owned by one user, and
// it is assigned with the unique identifier and (temporary) access key.
type Stream struct {
	Owner    string
	ID       string
	Name     string
	Json     string
	Type     int
	Metadata Metadata
}

// StreamsPage contains page related metadata as well as list of streams that
// belong to this page.
type StreamsPage struct {
	PageMetadata
	Streams []Stream
}

// StreamRepository specifies a stream persistence API.
type StreamRepository interface {
	// Save persists multiple streams. Streams are saved using a transaction. If one stream
	// fails then none will be saved. Successful operation is indicated by non-nil
	// error response.
	Save(context.Context, ...Stream) ([]Stream, error)

	// Update performs an update to the existing stream. A non-nil error is
	// returned to indicate operation failure.
	Update(context.Context, Stream) error

	// RetrieveByID retrieves the stream having the provided identifier, that is owned
	// by the specified user.
	RetrieveByID(context.Context, string, string) (Stream, error)

	// RetrieveAll retrieves the subset of streams owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (StreamsPage, error)

	// Remove removes the stream having the provided identifier, that is owned
	// by the specified user.
	Remove(context.Context, string, string) error
}

// StreamCache contains stream caching interface.
type StreamCache interface {
	// Save stores pair stream key, stream id.
	Save(context.Context, string, string) error

	// ID returns stream ID for given key.
	ID(context.Context, string) (string, error)

	// Removes stream from cache.
	Remove(context.Context, string) error
}
