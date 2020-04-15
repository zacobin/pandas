package http

import (
	"context"

	"github.com/cloustone/pandas/rulechain"
	"github.com/go-kit/kit/endpoint"
)

func addRuleChainEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRuleChainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		err := svc.AddNewRuleChain(ctx, req.token, req.rulechain)
		if err != nil {
			return nil, err
		}
		return updateRuleChainResponse{}, nil
	}
}

func rulechainInfoEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RuleChainRequestInfo)
		if err := req.validate(); err != nil {
			return nil, err
		}
		rulechain, err := svc.GetRuleChainInfo(ctx, req.token, req.RuleChainID)
		if err != nil {
			return nil, err
		}
		return rulechainResponse{rulechain}, nil
	}
}

func updateRuleChainEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRuleChainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		err := svc.UpdateRuleChain(ctx, req.token, req.rulechain)
		if err != nil {
			return nil, err
		}
		return updateRuleChainResponse{}, nil
	}
}

func deleteRuleChainEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RuleChainRequestInfo)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.RevokeRuleChain(ctx, req.token, req.RuleChainID)
		if err != nil {
			return nil, err
		}
		return updateRuleChainResponse{}, nil
	}
}

func listRuleChainEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRuleChainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		rulechains, err := svc.ListRuleChain(ctx, req.token)
		if err != nil {
			return nil, err
		}
		return listrulechainResponse{rulechains}, nil
	}
}

func startRuleChainEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RuleChainRequestInfo)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.StartRuleChain(ctx, req.token, req.RuleChainID)
		if err != nil {
			return nil, err
		}
		return updateRuleChainResponse{}, nil
	}
}

func stopRuleChainEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RuleChainRequestInfo)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.StopRuleChain(ctx, req.token, req.RuleChainID)
		if err != nil {
			return nil, err
		}
		return updateRuleChainResponse{}, nil
	}
}
