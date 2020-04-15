package http

import "github.com/cloustone/pandas/rulechain"

type RuleChainRequestInfo struct {
	UserID      string
	RuleChainID string
}

func (req RuleChainRequestInfo) validate() error {
	if req.UserID == "" {
		return rulechain.ErrUnauthorizedAccess
	}
	if req.RuleChainID == "" {
		return rulechain.ErrMalformedEntity
	}
	return nil
}

type updateRuleChainReq struct {
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
	UserID string
}

func (req listRuleChainReq) validate() error {
	if req.UserID == "" {
		return rulechain.ErrUnauthorizedAccess
	}
	return nil
}
