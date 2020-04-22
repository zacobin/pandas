// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"

	"github.com/cloustone/pandas/vms"
	"github.com/go-kit/kit/endpoint"
)

func addViewEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addViewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		view := vms.View{
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

func updateViewEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateViewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		view := vms.View{
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

func viewViewEndpoint(svc vms.Service) endpoint.Endpoint {
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

func listViewsEndpoint(svc vms.Service) endpoint.Endpoint {
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

func removeViewEndpoint(svc vms.Service) endpoint.Endpoint {
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

// Variable

func addVariableEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addVariableReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		variable := vms.Variable{
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

func updateVariableEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateVariableReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		variable := vms.Variable{
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

func viewVariableEndpoint(svc vms.Service) endpoint.Endpoint {
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

func listVariablesEndpoint(svc vms.Service) endpoint.Endpoint {
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

func removeVariableEndpoint(svc vms.Service) endpoint.Endpoint {
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

// Model

func addModelEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addModelReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		model := vms.Model{
			Name:     req.Name,
			Metadata: req.Metadata,
		}
		saved, err := svc.AddModel(ctx, req.token, model)
		if err != nil {
			return nil, err
		}

		res := modelRes{
			id:      saved.ID,
			created: true,
		}
		return res, nil
	}
}

func updateModelEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateModelReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		model := vms.Model{
			ID:       req.id,
			Name:     req.Name,
			Metadata: req.Metadata,
		}

		if err := svc.UpdateModel(ctx, req.token, model); err != nil {
			return nil, err
		}

		res := variableRes{id: req.id, created: false}
		return res, nil
	}
}

func viewModelEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewModelReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		view, err := svc.ViewModel(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := viewModelRes{
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

func listModelsEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listModelReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListModels(ctx, req.token, req.offset, req.limit, req.name, req.metadata)
		if err != nil {
			return nil, err
		}

		res := modelsPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
			Models: []viewModelRes{},
		}
		for _, variable := range page.Models {
			model := viewModelRes{
				Owner:    variable.Owner,
				ID:       variable.ID,
				Name:     variable.Name,
				Created:  variable.Created,
				Updated:  variable.Updated,
				Revision: variable.Revision,
				Metadata: variable.Metadata,
			}
			res.Models = append(res.Models, model)
		}

		return res, nil
	}
}

func removeModelEndpoint(svc vms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewModelReq)

		err := req.validate()
		if err != nil {
			return nil, err
		}

		if err := svc.RemoveModel(ctx, req.token, req.id); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}
