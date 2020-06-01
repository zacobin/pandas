// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/cloustone/pandas/alerts"
	"github.com/cloustone/pandas/mainflux"
)

var (
	_ mainflux.Response = (*listAlertsResponse)(nil)
	_ mainflux.Response = (*alertResponse)(nil)
	_ mainflux.Response = (*updateAlertResponse)(nil)
)

type genericResponse struct{}

func (res genericResponse) Code() int                  { return http.StatusOK }
func (res genericResponse) Headers() map[string]string { return map[string]string{} }
func (res genericResponse) Empty() bool                { return true }

type listAlertsResponse struct {
	Alerts []alerts.Alert `json:"alerts, omitempty"`
}

func (r listAlertsResponse) Code() int                  { return http.StatusOK }
func (r listAlertsResponse) Headers() map[string]string { return map[string]string{} }
func (r listAlertsResponse) Empty() bool                { return len(r.Alerts) > 0 }

type alertResponse struct {
	Alert alerts.Alert `json:"alert,omitempty"`
}

func (r alertResponse) Code() int                  { return http.StatusOK }
func (r alertResponse) Headers() map[string]string { return map[string]string{} }
func (r alertResponse) Empty() bool                { return r.Alert.Name == "" }

type updateAlertResponse struct{}

func (res updateAlertResponse) Code() int                  { return http.StatusOK }
func (res updateAlertResponse) Headers() map[string]string { return map[string]string{} }
func (res updateAlertResponse) Empty() bool                { return true }

type errorRes struct {
	Err string `json:"error"`
}

type principalAuthResponse struct{}

func (res principalAuthResponse) Code() int                  { return http.StatusOK }
func (res principalAuthResponse) Headers() map[string]string { return map[string]string{} }
func (res principalAuthResponse) Empty() bool                { return true }
