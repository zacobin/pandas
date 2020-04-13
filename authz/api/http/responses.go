// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package httpapi

import (
	"net/http"

	"github.com/cloustone/pandas/authz"
	"github.com/cloustone/pandas/mainflux"
)

var (
	_ mainflux.Response = (*listRolesResponse)(nil)
	_ mainflux.Response = (*roleResponse)(nil)
	_ mainflux.Response = (*updateRoleResponse)(nil)
)

type genericResponse struct{}

func (res genericResponse) Code() int                  { return http.StatusOK }
func (res genericResponse) Headers() map[string]string { return map[string]string{} }
func (res genericResponse) Empty() bool                { return true }

type listRolesResponse struct {
	Roles []authz.Role `json:"roles, omitempty"`
}

func (r listRolesResponse) Code() int                  { return http.StatusOK }
func (r listRolesResponse) Headers() map[string]string { return map[string]string{} }
func (r listRolesResponse) Empty() bool                { return len(r.Roles) > 0 }

type roleResponse struct {
	Role authz.Role `json:"role,omitempty"`
}

func (r roleResponse) Code() int                  { return http.StatusOK }
func (r roleResponse) Headers() map[string]string { return map[string]string{} }
func (r roleResponse) Empty() bool                { return false }

type updateRoleResponse struct{}

func (res updateRoleResponse) Code() int                  { return http.StatusOK }
func (res updateRoleResponse) Headers() map[string]string { return map[string]string{} }
func (res updateRoleResponse) Empty() bool                { return true }

type errorRes struct {
	Err string `json:"error"`
}
