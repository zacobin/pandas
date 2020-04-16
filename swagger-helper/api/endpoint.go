// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	swagger_helper "github.com/cloustone/pandas/swagger-helper"
	"github.com/go-kit/kit/endpoint"
)

func viewSwaggerEndpoint(svc swagger_helper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewSwaggerReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		swagger, err := svc.RetrieveDownstreamSwagger(ctx, req.token, req.module)
		if err != nil {
			return nil, err
		}

		instance := swagger.Host
		if !strings.HasPrefix(swagger.Host, "http") {
			instance = "http://" + swagger.Host
		}
		u, err := url.Parse(instance)
		if err != nil {
			return nil, err
		}
		u.Path = swagger.SwaggerUrl
		// TODO
		return nil, nil
	}
}

func listSwaggerEndpoint(svc swagger_helper.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listSwaggerReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		// TODO
		return nil, nil
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
