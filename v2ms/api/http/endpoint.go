// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"

	"github.com/cloustone/pandas/v2ms"
	"github.com/go-kit/kit/endpoint"
)

func addViewEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addViewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		view := v2ms.View{
			Name:     req.Name,
			Metadata: req.Metadata,
		}
		saved, err := svc.AddView(ctx, req.token, view)
		if err != nil {
			return nil, err
		}

		res := viewRes{
			id:      saved.ID,
			created: true,
		}
		return res, nil
	}
}

func updateViewEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateViewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		view := v2ms.View{
			ID:       req.id,
			Name:     req.Name,
			Metadata: req.Metadata,
		}

		if err := svc.UpdateView(ctx, req.token, view); err != nil {
			return nil, err
		}

		res := viewRes{id: req.id, created: false}
		return res, nil
	}
}

func viewViewEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewViewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		view, err := svc.ViewView(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := viewViewRes{
			Owner:    view.Owner,
			ID:       view.ID,
			Name:     view.Name,
			Created:  view.Created,
			Updated:  view.Updated,
			Revision: view.Revision,
			Metadata: view.Metadata,
		}
		return res, nil
	}
}

func listViewsEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listViewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListViews(ctx, req.token, req.offset, req.limit, req.name, req.metadata)
		if err != nil {
			return nil, err
		}

		res := viewsPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
			Views: []viewViewRes{},
		}
		for _, view := range page.Views {
			view := viewViewRes{
				Owner:    view.Owner,
				ID:       view.ID,
				Name:     view.Name,
				Created:  view.Created,
				Updated:  view.Updated,
				Revision: view.Revision,
				Metadata: view.Metadata,
			}
			res.Views = append(res.Views, view)
		}

		return res, nil
	}
}

func removeViewEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewViewReq)

		err := req.validate()
		if err != nil {
			return nil, err
		}

		if err := svc.RemoveView(ctx, req.token, req.id); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}

func addVariableEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addVariableReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		variable := v2ms.Variable{
			Name:     req.Name,
			Metadata: req.Metadata,
		}
		saved, err := svc.AddVariable(ctx, req.token, variable)
		if err != nil {
			return nil, err
		}

		res := variableRes{
			id:      saved.ID,
			created: true,
		}
		return res, nil
	}
}

func updateVariableEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateVariableReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		variable := v2ms.Variable{
			ID:       req.id,
			Name:     req.Name,
			Metadata: req.Metadata,
		}

		if err := svc.UpdateVariable(ctx, req.token, variable); err != nil {
			return nil, err
		}

		res := variableRes{id: req.id, created: false}
		return res, nil
	}
}

func viewVariableEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewVariableReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		view, err := svc.ViewVariable(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := viewVariableRes{
			Owner:    view.Owner,
			ID:       view.ID,
			Name:     view.Name,
			ThingID:  view.ThingID,
			Created:  view.Created,
			Updated:  view.Updated,
			Revision: view.Revision,
			Metadata: view.Metadata,
		}
		return res, nil
	}
}

func listVariablesEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listVariableReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListVariables(ctx, req.token, req.offset, req.limit, req.name, req.metadata)
		if err != nil {
			return nil, err
		}

		res := variablesPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
			Variables: []viewVariableRes{},
		}
		for _, variable := range page.Variables {
			variable := viewVariableRes{
				Owner:    variable.Owner,
				ID:       variable.ID,
				Name:     variable.Name,
				Created:  variable.Created,
				Updated:  variable.Updated,
				Revision: variable.Revision,
				Metadata: variable.Metadata,
			}
			res.Variables = append(res.Variables, variable)
		}

		return res, nil
	}
}

func removeVariableEndpoint(svc v2ms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewVariableReq)

		err := req.validate()
		if err != nil {
			return nil, err
		}

		if err := svc.RemoveView(ctx, req.token, req.id); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}
