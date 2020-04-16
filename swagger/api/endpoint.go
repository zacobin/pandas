// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	swagger "github.com/cloustone/pandas/swagger"
	"github.com/go-kit/kit/endpoint"
)

func viewSwaggerEndpoint(svc swagger.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewSwaggerReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		swagger, err := svc.RetrieveDownstreamSwagger(ctx, req.token, req.service)
		if err != nil {
			return nil, err
		}

		url := swagger.Host
		if !strings.HasPrefix(swagger.Host, "http") {
			url = "http://" + swagger.Host + swagger.SwaggerUrl
		}
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		return viewSwaggerRes{httpresp: resp}, nil
	}
}

func listSwaggerEndpoint(svc swagger.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listSwaggerReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		configs, err := svc.RetrieveSwaggerConfigs(context.TODO(), req.token)
		if err != nil {
			return nil, err
		}
		res := listSwaggerRes{
			Title:    configs.Info.Title,
			Version:  configs.Info.Version,
			Services: []string{},
		}
		for _, swagger := range configs.DownstreamSwaggers {
			res.Services = append(res.Services, swagger.Name)
		}
		return res, nil
	}
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
