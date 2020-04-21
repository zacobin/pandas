// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import "github.com/cloustone/pandas/vms"

const maxNameSize = 1024
const maxLimitSize = 100

type apiReq interface {
	validate() error
}

type addViewReq struct {
	token    string
	Name     string                 `json:"name,omitempty"`
	View     vms.View               `json:"view,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req addViewReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if len(req.Name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}

type updateViewReq struct {
	token    string
	id       string
	Name     string                 `json:"name,omitempty"`
	ThingID  string                 `json:"thing_id,omitempty"`
	View     vms.View               `json:"view,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req updateViewReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return vms.ErrMalformedEntity
	}

	if len(req.Name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}

type viewViewReq struct {
	token string
	id    string
}

func (req viewViewReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return vms.ErrMalformedEntity
	}

	return nil
}

type listViewReq struct {
	token    string
	offset   uint64
	limit    uint64
	name     string
	metadata map[string]interface{}
}

func (req *listViewReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.limit == 0 || req.limit > maxLimitSize {
		return vms.ErrMalformedEntity
	}

	if len(req.name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}

// Variable

type addVariableReq struct {
	token    string
	Name     string                 `json:"name,omitempty"`
	Variable vms.Variable           `json:"variable,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req addVariableReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if len(req.Name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}

type updateVariableReq struct {
	token    string
	id       string
	Name     string                 `json:"name,omitempty"`
	ThingID  string                 `json:"thing_id,omitempty"`
	Variable vms.Variable           `json:"view,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req updateVariableReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return vms.ErrMalformedEntity
	}

	if len(req.Name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}

type viewVariableReq struct {
	token string
	id    string
}

func (req viewVariableReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return vms.ErrMalformedEntity
	}

	return nil
}

type listVariableReq struct {
	token    string
	offset   uint64
	limit    uint64
	name     string
	metadata map[string]interface{}
}

func (req *listVariableReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.limit == 0 || req.limit > maxLimitSize {
		return vms.ErrMalformedEntity
	}

	if len(req.name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}

// Models

type addModelReq struct {
	token    string
	Name     string                 `json:"name,omitempty"`
	Model    vms.Model              `json:"variable,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req addModelReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if len(req.Name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}

type updateModelReq struct {
	token    string
	id       string
	Name     string                 `json:"name,omitempty"`
	Model    vms.Model              `json:"model,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req updateModelReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return vms.ErrMalformedEntity
	}

	if len(req.Name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}

type viewModelReq struct {
	token string
	id    string
}

func (req viewModelReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return vms.ErrMalformedEntity
	}

	return nil
}

type listModelReq struct {
	token    string
	offset   uint64
	limit    uint64
	name     string
	metadata map[string]interface{}
}

func (req *listModelReq) validate() error {
	if req.token == "" {
		return vms.ErrUnauthorizedAccess
	}

	if req.limit == 0 || req.limit > maxLimitSize {
		return vms.ErrMalformedEntity
	}

	if len(req.name) > maxNameSize {
		return vms.ErrMalformedEntity
	}

	return nil
}
