package kuiper

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/cloustone/pandas/kuiper/plugins"
	"github.com/cloustone/pandas/kuiper/xsql"
	"github.com/cloustone/pandas/kuiper/xsql/processors"
	"github.com/cloustone/pandas/mainflux"
)

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrScanMetadata indicates problem with metadata in db
	ErrScanMetadata = errors.New("Failed to scan metadata")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateStreams adds a list of streams to the user identified by the provided key.
	CreateStreams(context.Context, string, ...Stream) ([]Stream, error)

	// UpdateStream updates the stream identified by the provided ID, that
	// belongs to the user identified by the provided key.
	UpdateStream(context.Context, string, Stream) error

	// ViewStream retrieves data about the stream identified by the provided
	// ID, that belongs to the user identified by the provided key.
	ViewStream(context.Context, string, string) (Stream, error)

	// ListStreams retrieves data about subset of streams that belongs to the
	// user identified by the provided key.
	ListStreams(context.Context, string, uint64, uint64, string, Metadata) (StreamsPage, error)

	// RemoveStream removes the thing identified by the provided ID, that
	// belongs to the user identified by the provided key.
	RemoveStream(context.Context, string, string) error

	// CreateRules adds a list of things to the user identified by the provided key.
	CreateRules(context.Context, string, ...Rule) ([]Rule, error)

	// UpdateRule updates the thing identified by the provided ID, that
	// belongs to the user identified by the provided key.
	UpdateRule(context.Context, string, Rule) error

	// ViewRule retrieves data about the thing identified with the provided
	// ID, that belongs to the user identified by the provided key.
	ViewRule(context.Context, string, string) (Rule, error)

	// ListRules retrieves data about subset of things that belongs to the
	// user identified by the provided key.
	ListRules(context.Context, string, uint64, uint64, string, Metadata) (RulesPage, error)

	// RemoveRule removes the thing identified with the provided ID, that
	// belongs to the user identified by the provided key.
	RemoveRule(context.Context, string, string) error

	// StartRule start an already existed rule identifier with the provided ID,
	// that belongs to the user
	StartRule(context.Context, string, string) error

	// StopRule stop an already existed rule identifier with the provided ID,
	// that belongs to the user
	StopRule(context.Context, string, string) error

	// RestartRule restart an already existed rule identifier with the provided ID,
	// that belongs to the user
	RestartRule(context.Context, string, string) error

	// ListPluginSources retrieves data about subset of plugins sources
	ListPluginSources(context.Context, string, uint64, uint64, string, Metadata) (PluginSourcesPage, error)

	// ViewPluginSource retrieves data about the plugin source identified with the provided
	// ID.
	ViewPluginSource(context.Context, string, string) (PluginSource, error)

	// ListPluginSinks retrieves data about subset of plugins sinks
	ListPluginSinks(context.Context, string, uint64, uint64, string, Metadata) (PluginSinksPage, error)

	// ViewPluginSink retrieves data about the plugin sink identified with the provided
	// ID.
	ViewPluginSink(context.Context, string, string) (PluginSink, error)
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
	Name   string
}

var _ Service = (*kuiperService)(nil)

type kuiperService struct {
	auth            mainflux.AuthNServiceClient
	idp             IdentityProvider
	streams         StreamRepository
	rules           RuleRepository
	streamCache     StreamCache
	ruleCache       RuleCache
	ruleProcessor   *processors.RuleProcessor
	streamProcessor *processors.StreamProcessor
	pluginManager   *plugins.Manager
	ruleManager     *ruleManager
}

// New instantiates the things service implementation.
func New(auth mainflux.AuthNServiceClient, streams StreamRepository, rules RuleRepository,
	scache StreamCache, rcache RuleCache, idp IdentityProvider, pluginManager *plugins.Manager) Service {
	dataDir := "./"
	pluginManager, err := plugins.NewPluginManager()
	if err != nil {
		panic(err)
	}
	return &kuiperService{
		auth:            auth,
		idp:             idp,
		streams:         streams,
		rules:           rules,
		streamCache:     scache,
		ruleCache:       rcache,
		streamProcessor: processors.NewStreamProcessor(path.Join(path.Dir(dataDir), "stream")),
		pluginManager:   pluginManager,
		ruleManager:     newRuleManager(rules),
	}

}

// CreateStreams adds a list of streams to the user identified by the provided key.
func (ks *kuiperService) CreateStreams(ctx context.Context, token string, streams ...Stream) ([]Stream, error) {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return []Stream{}, ErrUnauthorizedAccess
	}
	for i := range streams {
		streams[i].ID, err = ks.idp.ID()
		if err != nil {
			return []Stream{}, err
		}
		streams[i].Owner = res.GetValue()
		parser := xsql.NewParser(strings.NewReader(streams[i].Json))
		stmt, err := xsql.Language.Parse(parser)
		if err != nil {
			return nil, err
		}
		switch stmt.(type) {
		case *xsql.StreamStmt:
		default:
			return nil, fmt.Errorf("Invalid stsream statement: %s", streams[i].Json)
		}
	}
	return ks.streams.Save(ctx, streams...)
}

// UpdateStream updates the stream identified by the provided ID, that
// belongs to the user identified by the provided key.
func (ks *kuiperService) UpdateStream(ctx context.Context, token string, stream Stream) error {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	stream.Owner = res.GetValue()
	return ks.streams.Update(ctx, stream)
}

// ViewStream retrieves data about the stream identified by the provided
// ID, that belongs to the user identified by the provided key.
func (ks *kuiperService) ViewStream(ctx context.Context, token string, id string) (Stream, error) {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Stream{}, ErrUnauthorizedAccess
	}
	stream, err := ks.streams.RetrieveByID(ctx, res.GetValue(), id)
	if err != nil {
		return Stream{}, err
	}
	parser := xsql.NewParser(strings.NewReader(stream.Json))
	stmt, err := xsql.Language.Parse(parser)
	if err != nil {
		return Stream{}, err
	}
	if _, ok := stmt.(*xsql.StreamStmt); !ok {
		return Stream{}, fmt.Errorf("Error resolving the stream %s, the data in db may be corrupted.", stream.Name)
	}
	return stream, nil
}

// ListStreams retrieves data about subset of streams that belongs to the
// user identified by the provided key.
func (ks *kuiperService) ListStreams(ctx context.Context, token string, offset, limit uint64, name string, metadata Metadata) (StreamsPage, error) {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return StreamsPage{}, ErrUnauthorizedAccess
	}
	streamsPage, err := ks.streams.RetrieveAll(ctx, res.GetValue(), offset, limit, name, metadata)
	if err != nil {
		return StreamsPage{}, err
	}
	for i := range streamsPage.Streams {
		s := streamsPage.Streams[i]
		parser := xsql.NewParser(strings.NewReader(s.Json))
		stmt, err := xsql.Language.Parse(parser)
		if err != nil {
			return StreamsPage{}, err
		}
		switch stmt.(type) {
		case *xsql.StreamStmt:
		default:
			return StreamsPage{}, fmt.Errorf("Invalid stsream statement: %s", s.Json)
		}

	}
	return streamsPage, nil
}

// RemoveStream removes the thing identified by the provided ID, that
// belongs to the user identified by the provided key.
func (ks *kuiperService) RemoveStream(ctx context.Context, token string, id string) error {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	return ks.streams.Remove(ctx, res.GetValue(), id)
}

// CreateRules adds a list of things to the user identified by the provided key.
func (ks *kuiperService) CreateRules(ctx context.Context, token string, rules ...Rule) ([]Rule, error) {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return []Rule{}, ErrUnauthorizedAccess
	}
	for i := range rules {
		rules[i].ID, err = ks.idp.ID()
		if err != nil {
			return []Rule{}, err
		}
		rules[i].Owner = res.GetValue()
	}

	// Save the created rules into repository at first
	if _, err := ks.rules.Save(ctx, rules...); err != nil {
		return []Rule{}, err
	}
	for i := range rules {
		r, err := ks.ruleManager.getRuleByJson(rules[i].Name, rules[i].SQL)
		if err != nil {
			return []Rule{}, err
		}
		//Start the rule
		rs, err := ks.ruleManager.createRuleState(r)
		if err != nil {
			return []Rule{}, err
		} else {
			err = ks.ruleManager.doStartRule(rs)
			if err != nil {
				return []Rule{}, err
			}
		}
	}
	return rules, nil
}

// UpdateRule updates the rule identified by the provided ID, that
// belongs to the user identified by the provided key.
func (ks *kuiperService) UpdateRule(ctx context.Context, token string, r Rule) error {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	r.Owner = res.GetValue()

	return ks.rules.Update(ctx, r)
}

// ViewRule retrieves data about the rule identified with the provided
// ID, that belongs to the user identified by the provided key.
func (ks *kuiperService) ViewRule(ctx context.Context, token string, id string) (Rule, error) {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Rule{}, ErrUnauthorizedAccess
	}

	rule, err := ks.rules.RetrieveByID(ctx, res.GetValue(), id)
	if err != nil {
		return Rule{}, err
	}
	if _, err = ks.ruleManager.getRuleByJson(rule.Name, rule.SQL); err != nil {
		return Rule{}, err
	}
	return rule, nil
}

// ListRules retrieves data about subset of things that belongs to the
// user identified by the provided key.
func (ks *kuiperService) ListRules(ctx context.Context, token string, offset, limit uint64, name string, metadata Metadata) (RulesPage, error) {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return RulesPage{}, ErrUnauthorizedAccess
	}

	return ks.rules.RetrieveAll(ctx, res.GetValue(), offset, limit, name, metadata)
}

// RemoveRule removes the thing identified with the provided ID, that
// belongs to the user identified by the provided key.
func (ks *kuiperService) RemoveRule(ctx context.Context, token string, id string) error {
	res, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	// Retrieve the rule with specified token and id
	r, err := ks.rules.RetrieveByID(ctx, token, id)
	if err != nil {
		return err
	}
	rule, err := ks.ruleManager.getRuleByJson(r.Name, r.SQL)
	if err != nil {
		return err
	}
	if err := ks.ruleManager.deleteRule(rule); err != nil {
		return err
	}
	return ks.rules.Remove(ctx, res.GetValue(), id)
}

// StartRule start an already existed rule identifier with the provided ID,
// that belongs to the user
func (ks *kuiperService) StartRule(ctx context.Context, token string, id string) error {
	_, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	// Retrieve the rule with specified token and id
	r, err := ks.rules.RetrieveByID(ctx, token, id)
	if err != nil {
		return err
	}
	rule, err := ks.ruleManager.getRuleByJson(r.Name, r.SQL)
	if err != nil {
		return err
	}
	return ks.ruleManager.startRule(rule)

}

// StopRule stop an already existed rule identifier with the provided ID,
// that belongs to the user
func (ks *kuiperService) StopRule(ctx context.Context, token string, id string) error {
	_, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	// Retrieve the rule with specified token and id
	r, err := ks.rules.RetrieveByID(ctx, token, id)
	if err != nil {
		return err
	}
	rule, err := ks.ruleManager.getRuleByJson(r.Name, r.SQL)
	if err != nil {
		return err
	}
	return ks.ruleManager.stopRule(rule)
}

// RestartRule restart an already existed rule identifier with the provided ID,
// that belongs to the user
func (ks *kuiperService) RestartRule(ctx context.Context, token string, id string) error {
	_, err := ks.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	// Retrieve the rule with specified token and id
	r, err := ks.rules.RetrieveByID(ctx, token, id)
	if err != nil {
		return err
	}
	rule, err := ks.ruleManager.getRuleByJson(r.Name, r.SQL)
	if err != nil {
		return err
	}
	return ks.ruleManager.restartRule(rule)
}

// ListPluginSources retrieves data about subset of plugins sources
func (ks *kuiperService) ListPluginSources(context.Context, string, uint64, uint64, string, Metadata) (PluginSourcesPage, error) {
	return PluginSourcesPage{}, nil
}

// ViewPluginSource retrieves data about the plugin source identified with the provided
// ID.
func (ks *kuiperService) ViewPluginSource(context.Context, string, string) (PluginSource, error) {
	return PluginSource{}, nil
}

// ListPluginSinks retrieves data about subset of plugins sinks
func (ks *kuiperService) ListPluginSinks(context.Context, string, uint64, uint64, string, Metadata) (PluginSinksPage, error) {
	return PluginSinksPage{}, nil
}

// ViewPluginSink retrieves data about the plugin sink identified with the provided
// ID.
func (ks *kuiperService) ViewPluginSink(context.Context, string, string) (PluginSink, error) {
	return PluginSink{}, nil
}
