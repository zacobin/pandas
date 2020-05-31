// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/realms"
)

var (
	_ mainflux.Response = (*listRealmsResponse)(nil)
	_ mainflux.Response = (*realmResponse)(nil)
	_ mainflux.Response = (*updateRealmResponse)(nil)
)

type genericResponse struct{}

func (res genericResponse) Code() int                  { return http.StatusOK }
func (res genericResponse) Headers() map[string]string { return map[string]string{} }
func (res genericResponse) Empty() bool                { return true }

type listRealmsResponse struct {
	Realms []realms.Realm `json:"realms, omitempty"`
}

func (r listRealmsResponse) Code() int                  { return http.StatusOK }
func (r listRealmsResponse) Headers() map[string]string { return map[string]string{} }
func (r listRealmsResponse) Empty() bool                { return len(r.Realms) > 0 }

type realmResponse struct {
	Realm realms.Realm `json:"realm,omitempty"`
}

func (r realmResponse) Code() int                  { return http.StatusOK }
func (r realmResponse) Headers() map[string]string { return map[string]string{} }
func (r realmResponse) Empty() bool                { return r.Realm.Name == "" }

type updateRealmResponse struct{}

func (res updateRealmResponse) Code() int                  { return http.StatusOK }
func (res updateRealmResponse) Headers() map[string]string { return map[string]string{} }
func (res updateRealmResponse) Empty() bool                { return true }

type errorRes struct {
	Err string `json:"error"`
}

type principalAuthResponse struct{}

func (res principalAuthResponse) Code() int                  { return http.StatusOK }
func (res principalAuthResponse) Headers() map[string]string { return map[string]string{} }
func (res principalAuthResponse) Empty() bool                { return true }
