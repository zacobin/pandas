// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"

	"github.com/cloustone/pandas/mainflux"
)

var (
	_ mainflux.Response = (*listCollectionsRes)(nil)
)

type listCollectionsRes struct {
	Products []string
}

func (res listCollectionsRes) Code() int {
	return http.StatusCreated
}

func (res listCollectionsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res listCollectionsRes) Empty() bool {
	return res.Products == nil
}
