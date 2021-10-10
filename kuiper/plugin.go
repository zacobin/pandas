package kuiper

// Stream represents a Mainflux thing. Each thing is owned by one user, and
// it is assigned with the unique identifier and (temporary) access key.
type Plugin struct {
	ID       string
	Name     string
	Json     string
	Type     int
	Stop     bool
	Metadata Metadata
}

type PluginSource struct {
	ID       string
	Name     string
	Json     string
	Type     int
	Stop     bool
	Metadata Metadata
}

// PluginSourcesPage contains page related metadata as well as list of plugin
// source that belong to this page.
type PluginSourcesPage struct {
	PageMetadata
	Sources []PluginSource
}

type PluginSink struct {
	ID       string
	Name     string
	Json     string
	Type     int
	Stop     bool
	Metadata Metadata
}

// PluginSinksPage contains page related metadata as well as list of plugin
// sink that belong to this page.
type PluginSinksPage struct {
	PageMetadata
	Sinks []PluginSink
}
