//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package lbs

import "context"

type LocationProvider interface {
	// Add monitored object's track points
	AddTrackPoint(context.Context, TrackPoint) error
	AddTrackPoints(context.Context, []TrackPoint) error

	// Create circle geofence and return goefence id if successful
	CreateCircleGeofence(context.Context, CircleGeofence) (string, error)

	// Update an existed geofence
	UpdateCircleGeofence(context.Context, CircleGeofence) error

	// Delete an existed geofence or monitored objects
	DeleteGeofence(context.Context, []string, []string) ([]string, error)

	// List geofences matched with ids or objects
	ListGeofence(context.Context, []string, []string) ([]*Geofence, error)

	// Add monitored object for specifed geofence
	AddMonitoredObject(context.Context, string, []string) error

	// Remove monitored object from specified geofence
	RemoveMonitoredObject(context.Context, string, []string) error

	// List monitored object in specifed geofence
	ListMonitoredObjects(ctx context.Context, fenceID string, pageIndex int, pageSize int) (int, []string)

	// Create poly geofence and return goefence id if successful
	CreatePolyGeofence(context.Context, PolyGeofence) (string, error)

	// Update an existed poly geofence
	UpdatePolyGeofence(context.Context, PolyGeofence) error

	// Alarms
	QueryStatus(ctx context.Context, monitoredPerson string, fenceIDs []string) (QueryStatus, error)
	GetHistoryAlarms(ctx context.Context, monitoredPerson string, fenceIDs []string) (HistoryAlarms, error)
	BatchGetHistoryAlarms(ctx context.Context, input *BatchGetHistoryAlarmsRequest) (BatchHistoryAlarmsResp, error)
	GetStayPoints(context.Context, *GetStayPointsRequest) (StayPoints, error)
	HandleAlarmNotification(context.Context, []byte) (*AlarmNotification, error)

	//Entity
	AddEntity(ctx context.Context, entityName string, entityDesc string) error
	UpdateEntity(ctx context.Context, entityName string, entityDesc string) error
	DeleteEntity(ctx context.Context, entityName string) error
	ListEntity(ctx context.Context, collectionID string, coordTypeOutput string, pageIndex int, pageSize int) (int, ListEntityStruct)
}

type LocationServingOptions struct {
	// Provider is location engine name, baidu or othere
	Provider string
	// AK is access key for lbs service provider
	AK string
	// ServiceID is service id for lbs service provider
	ServiceID string
}

func NewLocationServingOptions(provider, ak, serviceID string) LocationServingOptions {
	return LocationServingOptions{
		Provider:  provider,
		AK:        ak,
		ServiceID: serviceID,
	}
}
