// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cloustone/pandas/mainflux"
)

var (
	_ mainflux.Response = (*removeRes)(nil)
	_ mainflux.Response = (*projectRes)(nil)
	_ mainflux.Response = (*viewProjectRes)(nil)
	_ mainflux.Response = (*projectsPageRes)(nil)
)

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

type removeRes struct{}

func (res removeRes) Code() int {
	return http.StatusNoContent
}

func (res removeRes) Headers() map[string]string {
	return map[string]string{}
}

func (res removeRes) Empty() bool {
	return true
}

// Projects
type projectRes struct {
	id      string `json:"id"`
	created bool   `json:"created"`
}

func (res projectRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res projectRes) Headers() map[string]string {
	if res.created {
		return map[string]string{
			"Location": fmt.Sprintf("/views/%s", res.id),
		}
	}

	return map[string]string{}
}

func (res projectRes) Empty() bool {
	return true
}

type viewProjectRes struct {
	Owner    string                 `json:"owner,omitempty"`
	ID       string                 `json:"id"`
	ThingID  string                 `json:"thing_id"`
	Name     string                 `json:"name,omitempty"`
	Revision int                    `json:"revision"`
	Created  time.Time              `json:"created"`
	Updated  time.Time              `json:"updated"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (res viewProjectRes) Code() int {
	return http.StatusOK
}

func (res viewProjectRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewProjectRes) Empty() bool {
	return false
}

type projectsPageRes struct {
	pageRes
	Projects []viewProjectRes `json:"projects"`
}

func (res projectsPageRes) Code() int {
	return http.StatusOK
}

func (res projectsPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res projectsPageRes) Empty() bool { return false }
