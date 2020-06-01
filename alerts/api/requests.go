// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/cloustone/pandas/alerts"
)

const minPassLen = 8

type apiReq interface {
	validate() error
}

type createAlertReq struct {
	alert alerts.Alert
	token string
}

func (req createAlertReq) validate() error {
	return nil
}

type viewAlertInfoReq struct {
	token string
}

func (req viewAlertInfoReq) validate() error {
	if req.token == "" {
		return alerts.ErrUnauthorizedAccess
	}
	return nil
}

type alertRequestInfo struct {
	token     string
	alertName string
}

func (req alertRequestInfo) validate() error {
	if req.token == "" {
		return alerts.ErrUnauthorizedAccess
	}
	if req.alertName == "" {
		return alerts.ErrMalformedEntity
	}
	return nil
}

type updateAlertReq struct {
	token     string
	alertName string
	alert     alerts.Alert
}

func (req updateAlertReq) validate() error {
	if req.token == "" {
		return alerts.ErrUnauthorizedAccess
	}
	if req.alertName == "" {
		return alerts.ErrMalformedEntity
	}
	return nil
}

type principalAuthRequest struct {
	token string
}

func (req principalAuthRequest) validate() error {
	if req.token == "" {
		return alerts.ErrUnauthorizedAccess
	}
	return nil
}
