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
	_ mainflux.Response = (*viewRes)(nil)
	_ mainflux.Response = (*viewViewRes)(nil)
	_ mainflux.Response = (*viewsPageRes)(nil)
	_ mainflux.Response = (*removeRes)(nil)
	_ mainflux.Response = (*variableRes)(nil)
	_ mainflux.Response = (*viewVariableRes)(nil)
	_ mainflux.Response = (*variablesPageRes)(nil)
)

type viewRes struct {
	id      string
	created bool
}

func (res viewRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res viewRes) Headers() map[string]string {
	if res.created {
		return map[string]string{
			"Location": fmt.Sprintf("/views/%s", res.id),
		}
	}

	return map[string]string{}
}

func (res viewRes) Empty() bool {
	return true
}

type viewViewRes struct {
	Owner    string                 `json:"owner,omitempty"`
	ID       string                 `json:"id"`
	ThingID  string                 `json:"thing_id"`
	Name     string                 `json:"name,omitempty"`
	Revision int                    `json:"revision"`
	Created  time.Time              `json:"created"`
	Updated  time.Time              `json:"updated"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (res viewViewRes) Code() int {
	return http.StatusOK
}

func (res viewViewRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewViewRes) Empty() bool {
	return false
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

type viewsPageRes struct {
	pageRes
	Views []viewViewRes `json:"views"`
}

func (res viewsPageRes) Code() int {
	return http.StatusOK
}

func (res viewsPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewsPageRes) Empty() bool {
	return false
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

// Variables
type variableRes struct {
	id      string `json:"id"`
	created bool   `json:"created"`
}

func (res variableRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res variableRes) Headers() map[string]string {
	if res.created {
		return map[string]string{
			"Location": fmt.Sprintf("/views/%s", res.id),
		}
	}

	return map[string]string{}
}

func (res variableRes) Empty() bool {
	return true
}

type viewVariableRes struct {
	Owner    string                 `json:"owner,omitempty"`
	ID       string                 `json:"id"`
	ThingID  string                 `json:"thing_id"`
	Name     string                 `json:"name,omitempty"`
	Revision int                    `json:"revision"`
	Created  time.Time              `json:"created"`
	Updated  time.Time              `json:"updated"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (res viewVariableRes) Code() int {
	return http.StatusOK
}

func (res viewVariableRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewVariableRes) Empty() bool {
	return false
}

type variablesPageRes struct {
	pageRes
	Variables []viewVariableRes `json:"variables"`
}

func (res variablesPageRes) Code() int {
	return http.StatusOK
}

func (res variablesPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res variablesPageRes) Empty() bool { return false }
