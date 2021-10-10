package http

import (
	"context"

	"github.com/cloustone/pandas/kuiper"
	"github.com/go-kit/kit/endpoint"
)

func listStreamsEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listResourcesReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListStreams(ctx, req.token, req.offset, req.limit, req.name, req.metadata)
		if err != nil {
			return nil, err
		}

		res := streamsPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
			Streams: []viewStreamRes{},
		}
		for _, stream := range page.Streams {
			view := viewStreamRes{
				ID:       stream.ID,
				Name:     stream.Name,
				Json:     stream.Json,
				Metadata: stream.Metadata,
			}
			res.Streams = append(res.Streams, view)
		}

		return res, nil
	}
}

func createStreamEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		stream := kuiper.Stream{
			Name:     req.Name,
			Json:     req.Json,
			Metadata: req.Metadata,
		}
		saved, err := svc.CreateStreams(ctx, req.token, stream)
		if err != nil {
			return nil, err
		}

		res := streamRes{
			ID:      saved[0].ID,
			created: true,
		}
		return res, nil
	}
}

func viewStreamEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		stream, err := svc.ViewStream(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := viewStreamRes{
			ID:       stream.ID,
			Name:     stream.Name,
			Json:     stream.Json,
			Metadata: stream.Metadata,
		}
		return res, nil
	}
}

func deleteStreamEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		err := req.validate()
		if err == kuiper.ErrNotFound {
			return removeRes{}, nil
		}

		if err != nil {
			return nil, err
		}

		if err := svc.RemoveStream(ctx, req.token, req.id); err != nil {
			return nil, err
		}
		return removeRes{}, nil
	}
}

func listRulesEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listResourcesReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListRules(ctx, req.token, req.offset, req.limit, req.name, req.metadata)
		if err != nil {
			return nil, err
		}

		res := rulesPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
			Rules: []ruleRes{},
		}
		for _, r := range page.Rules {
			view := ruleRes{
				ID:       r.ID,
				Name:     r.Name,
				Metadata: r.Metadata,
			}
			res.Rules = append(res.Rules, view)
		}

		return res, nil
	}
}

func createRuleEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createRuleReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		stream := kuiper.Rule{
			Name:     req.Name,
			SQL:      req.SQL,
			Metadata: req.Metadata,
		}
		saved, err := svc.CreateRules(ctx, req.token, stream)
		if err != nil {
			return nil, err
		}

		res := streamRes{
			ID:      saved[0].ID,
			created: true,
		}
		return res, nil
	}
}

func viewRuleEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		r, err := svc.ViewRule(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := ruleRes{
			ID:       r.ID,
			Name:     r.Name,
			SQL:      r.SQL,
			Metadata: r.Metadata,
		}
		return res, nil
	}
}

func updateRuleEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRuleReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		r := kuiper.Rule{
			ID:       req.id,
			Name:     req.Name,
			SQL:      req.SQL,
			Metadata: req.Metadata,
		}
		err := svc.UpdateRule(ctx, req.token, r)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func deleteRuleEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		err := req.validate()
		if err == kuiper.ErrNotFound {
			return removeRes{}, nil
		}

		if err != nil {
			return nil, err
		}

		if err := svc.RemoveRule(ctx, req.token, req.id); err != nil {
			return nil, err
		}
		return removeRes{}, nil
	}
}

func viewRuleStatusEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		r, err := svc.ViewRule(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := ruleRes{
			ID:   r.ID,
			Name: r.Name,
			//Json:     r.Json,
			Metadata: r.Metadata,
		}
		return res, nil
	}
}

func startRuleEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ruleControlReq)
		var err error

		if err = req.validate(); err != nil {
			return nil, err
		}
		switch req.action {
		case ACTION_RULE_START:
			err = svc.StartRule(ctx, req.token, req.id)
		case ACTION_RULE_STOP:
			err = svc.StopRule(ctx, req.token, req.id)
		case ACTION_RULE_RESTART:
			err = svc.RestartRule(ctx, req.token, req.id)
		}

		return nil, err
	}
}

// Plugins
func listPluginSourcesEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func viewPluginSourceEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func listPluginSinksEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func viewPluginSinkEndpoint(svc kuiper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}
