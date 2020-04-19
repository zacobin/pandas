// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"github.com/cloustone/pandas/lbs"
)

type listCollectionsReq struct {
	token string
}

func (req listCollectionsReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	return nil
}

type CircleGeofence struct {
	Name             string
	MonitoredObjects []string
	Longitude        float64
	Latitude         float64
	Radius           float64
	CoordType        lbs.CoordType
	Denoise          int
	FenceID          string
}

type createCircleGeofenceReq struct {
	token     string
	projectID string
	fence     *CircleGeofence
}

func (req createCircleGeofenceReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if (req.projectID == "") || (req.fence == nil) {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type updateCircleGeofenceReq struct {
	token     string
	projectID string
	fence     *CircleGeofence
}

func (req updateCircleGeofenceReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if (req.projectID == "") || (req.fence == nil) {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type deleteGeofenceReq struct {
	token     string
	projectID string
	fenceIDs  []string
	objects   []string
}

func (req deleteGeofenceReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type listGeofencesReq struct {
	token     string
	projectID string
	fenceIDs  []string
	objects   []string
}

func (req listGeofencesReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type addMonitoredObjectReq struct {
	token     string
	projectID string
	fenceID   string
	objects   []string
}

func (req addMonitoredObjectReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type removeMonitoredObjectReq struct {
	token     string
	projectID string
	fenceID   string
	objects   []string
}

func (req removeMonitoredObjectReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type listMonitoredObjectsReq struct {
	token     string
	projectID string
	fenceID   string
	pageIndex int
	pageSize  int
}

func (req listMonitoredObjectsReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type PolyGeofence struct {
	Name             string
	MonitoredObjects []string
	Vertexes         string
	CoordType        lbs.CoordType
	Denoise          int
	FenceID          string
}

type createPolyGeofenceReq struct {
	token     string
	projectID string
	fence     *PolyGeofence
}

func (req createPolyGeofenceReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type updatePolyGeofenceReq struct {
	token     string
	projectID string
	fence     *PolyGeofence
}

func (req updatePolyGeofenceReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type getFenceIDsReq struct {
	token     string
	projectID string
}

func (req getFenceIDsReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type queryStatusReq struct {
	token           string
	projectID       string
	monitoredPerson string
	fenceIDs        []string
}

func (req queryStatusReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type getHistoryAlarmsReq struct {
	token           string
	projectID       string
	monitoredPerson string
	fenceIDs        []string
}

func (req getHistoryAlarmsReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type BatchGetHistoryAlarmsRequest struct {
	CoordTypeOutput string `protobuf:"bytes,3,opt,name=coord_type_output,json=coordTypeOutput" json:"coord_type_output,omitempty", bson:"coord_type_output,omitempty"`
	EndTime         string `protobuf:"bytes,4,opt,name=end_time,json=endTime" json:"end_time,omitempty", bson:"end_time,omitempty"`
	StartTime       string `protobuf:"bytes,5,opt,name=start_time,json=startTime" json:"start_time,omitempty", bson:"start_time,omitempty"`
	PageIndex       int  `protobuf:"varint,7,opt,name=page_index,json=pageIndex" json:"page_index,omitempty", bson:"page_index,omitempty"`
	PageSize        int  `protobuf:"varint,8,opt,name=page_size,json=pageSize" json:"page_size,omitempty", bson:"page_size,omitempty"`
}

type batchGetHistoryAlarmsReq struct {
	token     string
	projectID string
	input     *BatchGetHistoryAlarmsRequest
}

func (req batchGetHistoryAlarmsReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type GetStayPointsRequest struct {
	EndTime         string   `protobuf:"bytes,3,opt,name=end_time,json=endTime" json:"end_time,omitempty", bson:"end_time,omitempty"`
	EntityName      string   `protobuf:"bytes,4,opt,name=entity_name,json=entityName" json:"entity_name,omitempty", bson:"entity_name,omitempty"`
	FenceIDs        []string `protobuf:"bytes,5,rep,name=fence_ids,json=fenceIDs" json:"fence_ids,omitempty", bson:"fence_ids,omitempty"`
	PageIndex       int    `protobuf:"varint,6,opt,name=page_index,json=pageIndex" json:"page_index,omitempty", bson:"page_index,omitempty"`
	PageSize        int    `protobuf:"varint,7,opt,name=page_size,json=pageSize" json:"page_size,omitempty", bson:"page_size,omitempty"`
	StartTime       string   `protobuf:"bytes,8,opt,name=start_time,json=startTime" json:"start_time,omitempty", bson:"start_time,omitempty"`
	CoordTypeOutput string   `protobuf:"bytes,9,opt,name=coord_type_output,json=coordTypeOutput" json:"coord_type_output,omitempty", bson:"coord_type_output,omitempty"`
}

type getStayPointsReq struct {
	token     string
	projectID string
	input     *GetStayPointsRequest
}

func (req getStayPointsReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type notifyAlarmsReq struct {
	token     string
	projectID string
	content   []byte
}

func (req notifyAlarmsReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type getFenceUserIDReq struct {
	token     string
	projectID string
	fenceID   string
}

func (req getFenceUserIDReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type addEntityReq struct {
	token      string
	projectID  string
	entityName string
	entityDesc string
}

func (req addEntityReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type updateEntityReq struct {
	token      string
	projectID  string
	entityName string
	entityDesc string
}

func (req updateEntityReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type deleteEntityReq struct {
	token      string
	projectID  string
	entityName string
}

func (req deleteEntityReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}

type listEntityReq struct {
	token           string
	projectID       string
	coordTypeOutput string
	pageIndex       int
	pageSize        int
}

func (req listEntityReq) validate() error {
	if req.token == "" {
		return lbs.ErrUnauthorizedAccess
	}
	if req.projectID == "" {
		return lbs.ErrMalformedEntity
	}
	return nil
}
