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
		return addRuleChainResponse{}, nil
	}
}

func rulechainInfoEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RuleChainInfoRequest)
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
		rulechain, err := svc.UpdateRuleChain(ctx, req.token, req.rulechain)
		if err != nil {
			return nil, err
		}
		return updateRuleChainResponse{rulechain}, nil
	}
}

func deleteRuleChainEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RuleChainInfoRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.RevokeRuleChain(ctx, req.token, req.RuleChainID)
		if err != nil {
			return nil, err
		}
		return addRuleChainResponse{}, nil
	}
}

func listRuleChainEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRuleChainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListRuleChain(ctx, req.token, req.limit, req.offset)
		if err != nil {
			return nil, err
		}

		res := rulechainPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
			RuleChains: []rulechain.RuleChain{},
		}
		for _, rulechain := range page.RuleChains {
			res.RuleChains = append(res.RuleChains, rulechain)
		}

		return res, nil
	}
}

func updateRuleChainStatusEndpoint(svc rulechain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRuleChainStatusRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.UpdateRuleChainStatus(ctx, req.token, req.RuleChainID, req.updatestatus)
		if err != nil {
			return nil, err
		}
		return addRuleChainResponse{}, nil
	}
}
