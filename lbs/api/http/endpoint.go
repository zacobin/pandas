// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"

	"github.com/cloustone/pandas/lbs"
	"github.com/go-kit/kit/endpoint"
)

func listCollectionsEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listCollectionsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		saved, err := svc.ListCollections(ctx, req.token)
		if err != nil {
			return nil, err
		}

		res := listCollectionsRes{}
		for _, product := range saved {
			res.Products = append(res.Products, product)
		}
		return res, nil

	}
}
