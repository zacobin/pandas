// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/cloustone/pandas/realms"
)

const minPassLen = 8

type apiReq interface {
	validate() error
}

type createRealmReq struct {
	realm realms.Realm
	token string
}

func (req createRealmReq) validate() error {
	return req.realm.Validate()
}

type viewRealmInfoReq struct {
	token string
}

func (req viewRealmInfoReq) validate() error {
	if req.token == "" {
		return realms.ErrUnauthorizedAccess
	}
	return nil
}

type realmRequestInfo struct {
	token     string
	realmName string
}

func (req realmRequestInfo) validate() error {
	if req.token == "" {
		return realms.ErrUnauthorizedAccess
	}
	if req.realmName == "" {
		return realms.ErrMalformedEntity
	}
	return nil
}

type updateRealmReq struct {
	token     string
	realmName string
	realm     realms.Realm
}

func (req updateRealmReq) validate() error {
	if req.token == "" {
		return realms.ErrUnauthorizedAccess
	}
	if req.realmName == "" {
		return realms.ErrMalformedEntity
	}
	return nil
}

type principalAuthRequest struct {
	token     string
	principal realms.Principal
}

func (req principalAuthRequest) validate() error {
	if req.token == "" {
		return realms.ErrUnauthorizedAccess
	}
	if req.principal.Username == "" ||
		req.principal.Password == "" {
		return realms.ErrMalformedEntity
	}
	return nil
}
