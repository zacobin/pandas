// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"github.com/cloustone/pandas/alerts"
	"github.com/go-kit/kit/endpoint"
)

// Alert endpoint
func createAlertEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

func listAlertsEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

func alertInfoEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

func updateAlertEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

func deleteAlertEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

// Alerts rules
func createAlertRuleEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

func listAlertRulesEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

func alertRuleInfoEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

func updateAlertRuleEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}

func deleteAlertRuleEndpoint(svc alerts.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return genericResponse{}, nil
	}
}
