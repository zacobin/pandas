package http

import (
	"github.com/cloustone/pandas/kuiper"
)

const maxLimitSize = 100
const maxNameSize = 1024

type apiReq interface {
	validate() error
}

type createStreamReq struct {
	token    string
	Name     string                 `json:"name,noneempty"`
	Json     string                 `json:"json,noneempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req createStreamReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if len(req.Name) > maxNameSize {
		return kuiper.ErrMalformedEntity
	}

	return nil
}

type createStreamsReq struct {
	token   string
	Streams []createStreamReq
}

func (req createStreamsReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if len(req.Streams) <= 0 {
		return kuiper.ErrMalformedEntity
	}

	for _, thing := range req.Streams {
		if len(thing.Name) > maxNameSize {
			return kuiper.ErrMalformedEntity
		}
	}

	return nil
}

type updateStreamReq struct {
	token    string
	id       string
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req updateStreamReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return kuiper.ErrMalformedEntity
	}

	if len(req.Name) > maxNameSize {
		return kuiper.ErrMalformedEntity
	}

	return nil
}

type createRuleReq struct {
	token    string
	Name     string                 `json:"name,omitempty"`
	SQL      string                 `json:"sql,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req createRuleReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if len(req.Name) > maxNameSize {
		return kuiper.ErrMalformedEntity
	}

	return nil
}

type createRulesReq struct {
	token string
	Rules []createRuleReq
}

func (req createRulesReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if len(req.Rules) <= 0 {
		return kuiper.ErrMalformedEntity
	}

	for _, rule := range req.Rules {
		if len(rule.Name) > maxNameSize {
			return kuiper.ErrMalformedEntity
		}
	}

	return nil
}

type updateRuleReq struct {
	token    string
	id       string
	Name     string                 `json:"name,omitempty"`
	SQL      string                 `json:"sql,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req updateRuleReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return kuiper.ErrMalformedEntity
	}

	if len(req.Name) > maxNameSize {
		return kuiper.ErrMalformedEntity
	}

	return nil
}

type viewResourceReq struct {
	token string
	id    string
}

func (req viewResourceReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return kuiper.ErrMalformedEntity
	}

	return nil
}

type listResourcesReq struct {
	token    string
	offset   uint64
	limit    uint64
	name     string
	metadata map[string]interface{}
}

func (req *listResourcesReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if req.limit == 0 || req.limit > maxLimitSize {
		return kuiper.ErrMalformedEntity
	}

	if len(req.name) > maxNameSize {
		return kuiper.ErrMalformedEntity
	}

	return nil
}

type ruleStatusReq struct {
	token string
	id    string
}

func (req ruleStatusReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return kuiper.ErrMalformedEntity
	}

	return nil
}

const (
	ACTION_RULE_START   = "start"
	ACTION_RULE_STOP    = "stop"
	ACTION_RULE_RESTART = "restart"
)

type ruleControlReq struct {
	token  string
	id     string
	action string
}

func (req ruleControlReq) validate() error {
	if req.token == "" {
		return kuiper.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return kuiper.ErrMalformedEntity
	}
	switch req.action {
	case ACTION_RULE_START:
	case ACTION_RULE_STOP:
	case ACTION_RULE_RESTART:
	default:
		return kuiper.ErrMalformedEntity
	}

	return nil
}
