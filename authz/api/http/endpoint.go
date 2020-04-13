// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package httpapi

import (
	"context"

	"github.com/cloustone/pandas/authz"
	"github.com/go-kit/kit/endpoint"
)

// listRolesEndpoint return all roles in authz
func rolesInfoEndpoint(svc authz.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(genericRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		roles, err := svc.ListRoles(ctx, req.token)
		if err != nil {
			return listRolesResponse{}, err
		}
		return listRolesResponse{Roles: roles}, nil
	}
}

func roleInfoEndpoint(svc authz.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewRoleInfoRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}
		role, err := svc.RetrieveRole(ctx, req.token, req.roleName)
		if err != nil {
			return nil, err
		}
		return roleResponse{Role: role}, nil
	}
}

func updateRoleEndpoint(svc authz.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRoleRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.UpdateRole(ctx, req.token, req.role)
		if err != nil {
			return nil, err
		}
		return genericResponse{}, nil
	}
}

func authzEndpoint(svc authz.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(authzRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}
		err := svc.Authorize(ctx, req.token, req.roleName, authz.Subject{Object: req.subject})
		if err != nil {
			return nil, err
		}
		return genericResponse{}, nil
	}
}
