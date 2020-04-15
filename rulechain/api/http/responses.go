package http

import (
	"net/http"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/rulechain"
)

var (
	_ mainflux.Response = (*updateRuleChainResponse)(nil)
)

type updateRuleChainResponse struct{}

func (res updateRuleChainResponse) Code() int                  { return http.StatusOK }
func (res updateRuleChainResponse) Headers() map[string]string { return map[string]string{} }
func (res updateRuleChainResponse) Empty() bool                { return true }

type rulechainResponse struct {
	RuleChain rulechain.RuleChain `json:"rulechain,omitempty`
}

func (r rulechainResponse) Code() int                  { return http.StatusOK }
func (r rulechainResponse) Headers() map[string]string { return map[string]string{} }
func (r rulechainResponse) Empty() bool                { return r.RuleChain.ID == "" }

type listrulechainResponse struct {
	RuleChains []rulechain.RuleChain `json:"rulechains,omitempty`
}

func (r listrulechainResponse) Code() int                  { return http.StatusOK }
func (r listrulechainResponse) Headers() map[string]string { return map[string]string{} }
func (r listrulechainResponse) Empty() bool                { return len(r.RuleChains) > 0 }

type errorRes struct {
	Err string `json:"error"`
}
