package http

import (
	"github.com/cloustone/pandas/rulechain"
)

type RuleChainInfoRequest struct {
	token       string
	RuleChainID string
}

func (req RuleChainInfoRequest) validate() error {
	if req.token == "" {
		return rulechain.ErrUnauthorizedAccess
	}
	if req.RuleChainID == "" {
		return rulechain.ErrMalformedEntity
	}
	return nil
}

type updateRuleChainStatusRequest struct {
	token        string
	RuleChainID  string
	updatestatus string
}

func (req updateRuleChainStatusRequest) validate() error {
	if req.token == "" {
		return rulechain.ErrUnauthorizedAccess
	}
	if req.RuleChainID == "" {
		return rulechain.ErrMalformedEntity
	}
	return nil
}

type updateRuleChainReq struct {
	token     string
	rulechain rulechain.RuleChain
}

func (req updateRuleChainReq) validate() error {
	if req.rulechain.UserID == "" {
		return rulechain.ErrUnauthorizedAccess
	}
	if req.rulechain.ID == "" {
		return rulechain.ErrMalformedEntity
	}
	return nil
}

type listRuleChainReq struct {
	token  string
	offset uint64
	limit  uint64
}

func (req listRuleChainReq) validate() error {
	if req.token == "" {
		return rulechain.ErrUnauthorizedAccess
	}
	return nil
}
