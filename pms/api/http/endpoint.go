// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"

	"github.com/cloustone/pandas/pms"
	"github.com/go-kit/kit/endpoint"
)

func addProjectEndpoint(svc pms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addProjectReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		project := pms.Project{
			Name:     req.Name,
			Metadata: req.Metadata,
		}
		saved, err := svc.AddProject(ctx, req.token, project)
		if err != nil {
			return nil, err
		}

		res := projectRes{
			id:      saved.ID,
			created: true,
		}
		return res, nil
	}
}

func updateProjectEndpoint(svc pms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateProjectReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		project := pms.Project{
			ID:       req.id,
			Name:     req.Name,
			Metadata: req.Metadata,
		}

		if err := svc.UpdateProject(ctx, req.token, project); err != nil {
			return nil, err
		}

		res := projectRes{id: req.id, created: false}
		return res, nil
	}
}

func viewProjectEndpoint(svc pms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewProjectReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		view, err := svc.ViewProject(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := viewProjectRes{
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

func listProjectsEndpoint(svc pms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listProjectReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListProjects(ctx, req.token, req.offset, req.limit, req.name, req.metadata)
		if err != nil {
			return nil, err
		}

		res := projectsPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
			Projects: []viewProjectRes{},
		}
		for _, project := range page.Projects {
			project := viewProjectRes{
				Owner:    project.Owner,
				ID:       project.ID,
				Name:     project.Name,
				Created:  project.Created,
				Updated:  project.Updated,
				Revision: project.Revision,
				Metadata: project.Metadata,
			}
			res.Projects = append(res.Projects, project)
		}

		return res, nil
	}
}

func removeProjectEndpoint(svc pms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewProjectReq)

		err := req.validate()
		if err != nil {
			return nil, err
		}

		if err := svc.RemoveProject(ctx, req.token, req.id); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}
