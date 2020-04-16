// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/cloustone/pandas/mainflux"
)

var (
	_ mainflux.Response = (*viewSwaggerRes)(nil)
)

type viewSwaggerRes struct {
	httpresp *http.Response
}

func (res viewSwaggerRes) Code() int {
	return http.StatusOK
}

func (res viewSwaggerRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewSwaggerRes) Empty() bool {
	return true
}

type listSwaggerRes struct {
	Title    string   `json:"title"`
	Version  string   `json:"version"`
	Services []string `json:"services"`
}

func (res listSwaggerRes) Code() int {
	return http.StatusOK
}

func (res listSwaggerRes) Headers() map[string]string {
	return map[string]string{}
}

func (res listSwaggerRes) Empty() bool {
	return true
}
