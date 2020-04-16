// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	swagger "github.com/cloustone/pandas/swagger"
)

const maxNameSize = 1024
const maxLimitSize = 100

type apiReq interface {
	validate() error
}

type listSwaggerReq struct {
	token string
}

func (req listSwaggerReq) validate() error {
	if req.token == "" {
		return swagger.ErrUnauthorizedAccess
	}

	return nil
}

type viewSwaggerReq struct {
	token   string
	service string
	httpreq *http.Request
}

func (req viewSwaggerReq) validate() error {
	if req.token == "" {
		return swagger.ErrUnauthorizedAccess
	}

	if req.service == "" {
		return swagger.ErrMalformedEntity
	}

	return nil
}
