// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	swagger_helper "github.com/cloustone/pandas/swagger-helper"
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
		return swagger_helper.ErrUnauthorizedAccess
	}

	return nil
}

type viewSwaggerReq struct {
	token  string
	module string
}

func (req viewSwaggerReq) validate() error {
	if req.token == "" {
		return swagger_helper.ErrUnauthorizedAccess
	}

	if req.module == "" {
		return swagger_helper.ErrMalformedEntity
	}

	return nil
}
