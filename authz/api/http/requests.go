// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package httpapi

import "github.com/cloustone/pandas/authz"

const minPassLen = 8

type apiReq interface {
	validate() error
}

type genericRequest struct {
	token string
}

func (req genericRequest) validate() error {
	if req.token == "" {
		return authz.ErrUnauthorizedAccess
	}
	return nil
}

type viewRoleInfoRequest struct {
	token    string
	roleName string
}

func (req viewRoleInfoRequest) validate() error {
	if req.token == "" {
		return authz.ErrUnauthorizedAccess
	}
	if req.roleName == "" {
		return authz.ErrMalformedEntity
	}
	return nil
}

type updateRoleRequest struct {
	token string
	role  authz.Role
}

func (req updateRoleRequest) validate() error {
	if req.token == "" {
		return authz.ErrUnauthorizedAccess
	}
	return req.role.Validate()
}

type authzRequest struct {
	token    string `json:'-'`
	roleName string `json:"roleName"`
	subject  string `json:"subject"`
}

func (req authzRequest) validate() error {
	if req.token == "" {
		return authz.ErrUnauthorizedAccess
	}
	if req.roleName == "" || req.subject == "" {
		return authz.ErrMalformedEntity
	}
	return nil
}
