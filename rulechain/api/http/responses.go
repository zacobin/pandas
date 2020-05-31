package http

import (
	"net/http"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/rulechain"
)

var (
	_ mainflux.Response = (*updateRuleChainResponse)(nil)
)

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

type addRuleChainResponse struct{}

func (res addRuleChainResponse) Code() int                  { return http.StatusOK }
func (res addRuleChainResponse) Headers() map[string]string { return map[string]string{} }
func (res addRuleChainResponse) Empty() bool                { return true }

type updateRuleChainResponse struct {
	RuleChain rulechain.RuleChain `json:"rulechain,omitempty`
}

func (res updateRuleChainResponse) Code() int                  { return http.StatusOK }
func (res updateRuleChainResponse) Headers() map[string]string { return map[string]string{} }
func (res updateRuleChainResponse) Empty() bool                { return true }

type rulechainResponse struct {
	RuleChain rulechain.RuleChain `json:"rulechain,omitempty`
}

func (r rulechainResponse) Code() int                  { return http.StatusOK }
func (r rulechainResponse) Headers() map[string]string { return map[string]string{} }
func (r rulechainResponse) Empty() bool                { return r.RuleChain.ID == "" }

type rulechainPageRes struct {
	pageRes
	RuleChains []rulechain.RuleChain
}

func (r rulechainPageRes) Code() int                  { return http.StatusOK }
func (r rulechainPageRes) Headers() map[string]string { return map[string]string{} }
func (r rulechainPageRes) Empty() bool                { return true }

type errorRes struct {
	Err string `json:"error"`
}
