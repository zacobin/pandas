// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"github.com/cloustone/pandas/realms"
	"github.com/go-kit/kit/endpoint"
)

func registrationEndpoint(svc realms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createRealmReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.Register(ctx, req.realm); err != nil {
			return genericResponse{}, err
		}
		return genericResponse{}, nil
	}
}

func listRealmEndpoint(svc realms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewRealmInfoReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		realms, err := svc.ListRealms(ctx, req.token)
		if err != nil {
			return nil, err
		}
		return listRealmsResponse{Realms: realms}, nil
	}
}

func realmInfoEndpoint(svc realms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(realmRequestInfo)
		if err := req.validate(); err != nil {
			return nil, err
		}

		realm, err := svc.RealmInfo(ctx, req.token, req.realmName)
		if err != nil {
			return nil, err
		}
		return realmResponse{realm}, nil
	}
}

func updateRealmEndpoint(svc realms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRealmReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		err := svc.UpdateRealm(ctx, req.token, req.realm)
		if err != nil {
			return nil, err
		}
		return updateRealmResponse{}, nil
	}
}

func deleteRealmEndpoint(svc realms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(realmRequestInfo)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.RevokeRealm(ctx, req.token, req.realmName)
		if err != nil {
			return nil, err
		}
		return updateRealmResponse{}, nil
	}
}

func principalAuthEndpoint(svc realms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(principalAuthRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.Identify(ctx, req.token, req.principal)
		if err != nil {
			return nil, err
		}
		return principalAuthResponse{}, nil
	}
}
