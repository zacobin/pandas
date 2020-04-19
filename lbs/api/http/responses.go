// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"

	"github.com/cloustone/pandas/lbs"
	"github.com/cloustone/pandas/mainflux"
)

var (
	_ mainflux.Response = (*listCollectionsRes)(nil)
	_ mainflux.Response = (*createCircleGeofenceRes)(nil)
	_ mainflux.Response = (*updateCircleGeofenceRes)(nil)
	_ mainflux.Response = (*deleteGeofenceRes)(nil)
	_ mainflux.Response = (*listGeofencesRes)(nil)
	_ mainflux.Response = (*addMonitoredObjectRes)(nil)
	_ mainflux.Response = (*removeMonitoredObjectRes)(nil)
	_ mainflux.Response = (*listMonitoredObjectsRes)(nil)
	_ mainflux.Response = (*createPolyGeofenceRes)(nil)
	_ mainflux.Response = (*updatePolyGeofenceRes)(nil)
	_ mainflux.Response = (*getFenceIDsRes)(nil)
	_ mainflux.Response = (*queryStatusRes)(nil)
	_ mainflux.Response = (*getHistoryAlarmsRes)(nil)
	_ mainflux.Response = (*batchGetHistoryAlarmsRes)(nil)
	_ mainflux.Response = (*getStayPointsRes)(nil)
	_ mainflux.Response = (*notifyAlarmsRes)(nil)
	_ mainflux.Response = (*getFenceUserIDRes)(nil)
	_ mainflux.Response = (*addEntityRes)(nil)
	_ mainflux.Response = (*updateEntityRes)(nil)
	_ mainflux.Response = (*deleteEntityRes)(nil)
	_ mainflux.Response = (*listEntityRes)(nil)
)

type listCollectionsRes struct {
	Products []string
}

func (res listCollectionsRes) Code() int {
	return http.StatusCreated
}

func (res listCollectionsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res listCollectionsRes) Empty() bool {
	return res.Products == nil
}

type createCircleGeofenceRes struct {
	fenceID string
}

func (res createCircleGeofenceRes) Code() int {
	return http.StatusCreated
}

func (res createCircleGeofenceRes) Headers() map[string]string {
	return map[string]string{}
}

func (res createCircleGeofenceRes) Empty() bool {
	return res.fenceID == ""
}

type updateCircleGeofenceRes struct {
	updated bool
}

func (res updateCircleGeofenceRes) Code() int {
	return http.StatusCreated
}

func (res updateCircleGeofenceRes) Headers() map[string]string {
	return map[string]string{}
}

func (res updateCircleGeofenceRes) Empty() bool {
	return false
}

type deleteGeofenceRes struct {
	deleted bool
}

func (res deleteGeofenceRes) Code() int {
	return http.StatusCreated
}

func (res deleteGeofenceRes) Headers() map[string]string {
	return map[string]string{}
}

func (res deleteGeofenceRes) Empty() bool {
	return false
}

type Vertexe struct {
	Longitude float64
	Latitude  float64
}

type Geofence struct {
	FenceID         string
	FenceName       string
	MonitoredObject []string
	Shape           string
	Longitude       float64
	Latitude        float64
	Radius          float64
	CoordType       lbs.CoordType
	Denoise         int
	CreateTime      string
	UpdateTime      string
	Vertexes        []*Vertexe
}

type listGeofencesRes struct {
	fenceList []*Geofence
}

func (res listGeofencesRes) Code() int {
	return http.StatusCreated
}

func (res listGeofencesRes) Headers() map[string]string {
	return map[string]string{}
}

func (res listGeofencesRes) Empty() bool {
	return res.fenceList == nil
}

type addMonitoredObjectRes struct{}

func (res addMonitoredObjectRes) Code() int {
	return http.StatusCreated
}

func (res addMonitoredObjectRes) Headers() map[string]string {
	return map[string]string{}
}

func (res addMonitoredObjectRes) Empty() bool {
	return true
}

type removeMonitoredObjectRes struct{}

func (res removeMonitoredObjectRes) Code() int {
	return http.StatusCreated
}

func (res removeMonitoredObjectRes) Headers() map[string]string {
	return map[string]string{}
}

func (res removeMonitoredObjectRes) Empty() bool {
	return true
}

type listMonitoredObjectsRes struct {
	total   int
	objects []string
}

func (res listMonitoredObjectsRes) Code() int {
	return http.StatusCreated
}

func (res listMonitoredObjectsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res listMonitoredObjectsRes) Empty() bool {
	return res.objects == nil
}

type createPolyGeofenceRes struct {
	fenceID string
}

func (res createPolyGeofenceRes) Code() int {
	return http.StatusCreated
}

func (res createPolyGeofenceRes) Headers() map[string]string {
	return map[string]string{}
}

func (res createPolyGeofenceRes) Empty() bool {
	return res.fenceID == ""
}

type updatePolyGeofenceRes struct{}

func (res updatePolyGeofenceRes) Code() int {
	return http.StatusCreated
}

func (res updatePolyGeofenceRes) Headers() map[string]string {
	return map[string]string{}
}

func (res updatePolyGeofenceRes) Empty() bool {
	return true
}

type getFenceIDsRes struct {
	fenceIDs []string
}

func (res getFenceIDsRes) Code() int {
	return http.StatusCreated
}

func (res getFenceIDsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res getFenceIDsRes) Empty() bool {
	return res.fenceIDs == nil
}

type MonitoredStatus struct {
	FenceID         int    `json:"fence_id"`
	MonitoredStatus string `json:"monitored_status"`
}

type queryStatusRes struct {
	Status            int
	Message           string
	Size              int
	MonitoredStatuses []MonitoredStatus
}

func (res queryStatusRes) Code() int {
	return http.StatusCreated
}

func (res queryStatusRes) Headers() map[string]string {
	return map[string]string{}
}

func (res queryStatusRes) Empty() bool {
	return false
}

type Alarm struct {
	FenceID          string     `json:"fence_id,noempty"`
	FenceName        string     `json:"fence_name,noempty"`
	MonitoredObjects []string   `json:"monitored_objexts"`
	Action           string     `json:"action"`
	AlarmPoint       AlarmPoint `json:"alarm_point"`
	PrePoint         AlarmPoint `json:"pre_point"`
}
type AlarmPoint struct {
	Longitude  float64       `json:"longitude"`
	Latitude   float64       `json:"latitude"`
	Radius     int           `json:"radius"`
	CoordType  lbs.CoordType `json:"coord_type"`
	LocTime    string        `json:"loc_time"`
	CreateTime string        `json:"create_time"`
}

type getHistoryAlarmsRes struct {
	Status  int
	Message string
	Total   int
	Size    int
	Alarms  []*Alarm
}

func (res getHistoryAlarmsRes) Code() int {
	return http.StatusCreated
}

func (res getHistoryAlarmsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res getHistoryAlarmsRes) Empty() bool {
	return false
}

type batchGetHistoryAlarmsRes struct {
	Status  int
	Message string
	Total   int
	Size    int
	Alarms  []*Alarm
}

func (res batchGetHistoryAlarmsRes) Code() int {
	return http.StatusCreated
}

func (res batchGetHistoryAlarmsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res batchGetHistoryAlarmsRes) Empty() bool {
	return false
}

type Point struct {
	Longitude float64       `json:"longitude"`
	Latitude  float64       `json:"latitude"`
	CoordType lbs.CoordType `json:"coord_type"`
	LocTime   string        `json:"loc_time"`
}

type getStayPointsRes struct {
	Status     int
	Message    string
	Total      int
	Size       int
	Distance   int
	EndPoint   *Point
	StartPoint *Point
	Points     []*Point
}

func (res getStayPointsRes) Code() int {
	return http.StatusCreated
}

func (res getStayPointsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res getStayPointsRes) Empty() bool {
	return false
}

type notifyAlarmsRes struct{}

func (res notifyAlarmsRes) Code() int {
	return http.StatusCreated
}

func (res notifyAlarmsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res notifyAlarmsRes) Empty() bool {
	return true
}

type getFenceUserIDRes struct {
	UserID string
}

func (res getFenceUserIDRes) Code() int {
	return http.StatusCreated
}

func (res getFenceUserIDRes) Headers() map[string]string {
	return map[string]string{}
}

func (res getFenceUserIDRes) Empty() bool {
	return res.UserID == ""
}

type addEntityRes struct {
	Successed bool
}

func (res addEntityRes) Code() int {
	return http.StatusCreated
}

func (res addEntityRes) Headers() map[string]string {
	return map[string]string{}
}

func (res addEntityRes) Empty() bool {
	return false
}

type deleteEntityRes struct {
	Successed bool
}

func (res deleteEntityRes) Code() int {
	return http.StatusCreated
}

func (res deleteEntityRes) Headers() map[string]string {
	return map[string]string{}
}

func (res deleteEntityRes) Empty() bool {
	return false
}

type updateEntityRes struct {
	Successed bool
}

func (res updateEntityRes) Code() int {
	return http.StatusCreated
}

func (res updateEntityRes) Headers() map[string]string {
	return map[string]string{}
}

func (res updateEntityRes) Empty() bool {
	return false
}

type EntityInfo struct {
	EntityName string
	Latitude   float64
	Longitude  float64
}

type listEntityRes struct {
	Total       int
	EntityInfos []*EntityInfo
}

func (res listEntityRes) Code() int {
	return http.StatusCreated
}

func (res listEntityRes) Headers() map[string]string {
	return map[string]string{}
}

func (res listEntityRes) Empty() bool {
	return false
}
