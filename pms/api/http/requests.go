// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import "github.com/cloustone/pandas/pms"

const maxNameSize = 1024
const maxLimitSize = 100

type apiReq interface {
	validate() error
}

// Project

type addProjectReq struct {
	token    string
	Name     string                 `json:"name,omitempty"`
	Project  pms.Project            `json:"project,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req addProjectReq) validate() error {
	if req.token == "" {
		return pms.ErrUnauthorizedAccess
	}

	if len(req.Name) > maxNameSize {
		return pms.ErrMalformedEntity
	}

	return nil
}

type updateProjectReq struct {
	token    string
	id       string
	Name     string                 `json:"name,omitempty"`
	ThingID  string                 `json:"thing_id,omitempty"`
	Project  pms.Project            `json:"view,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req updateProjectReq) validate() error {
	if req.token == "" {
		return pms.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return pms.ErrMalformedEntity
	}

	if len(req.Name) > maxNameSize {
		return pms.ErrMalformedEntity
	}

	return nil
}

type viewProjectReq struct {
	token string
	id    string
}

func (req viewProjectReq) validate() error {
	if req.token == "" {
		return pms.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return pms.ErrMalformedEntity
	}

	return nil
}

type listProjectReq struct {
	token    string
	offset   uint64
	limit    uint64
	name     string
	metadata map[string]interface{}
}

func (req *listProjectReq) validate() error {
	if req.token == "" {
		return pms.ErrUnauthorizedAccess
	}

	if req.limit == 0 || req.limit > maxLimitSize {
		return pms.ErrMalformedEntity
	}

	if len(req.name) > maxNameSize {
		return pms.ErrMalformedEntity
	}

	return nil
}
