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

type LocationProvider interface {
	// Add monitored object's track points
	AddTrackPoint(point TrackPoint)
	AddTrackPoints(points []TrackPoint)

	// Create circle geofence and return goefence id if successful
	CreateCircleGeofence(c CircleGeofence) (string, error)

	// Update an existed geofence
	UpdateCircleGeofence(c CircleGeofence) error

	// Delete an existed geofence or monitored objects
	DeleteGeofence(fenceIds []string, objects []string) ([]string, error)

	// List geofences matched with ids or objects
	ListGeofence(fenceIds []string, objects []string) ([]*Geofence, error)

	// Add monitored object for specifed geofence
	AddMonitoredObject(fenceId string, objects []string) error

	// Remove monitored object from specified geofence
	RemoveMonitoredObject(fenceId string, objects []string) error

	// List monitored object in specifed geofence
	ListMonitoredObjects(fenceId string, pageIndex int, pageSize int) (int, []string)

	// Create poly geofence and return goefence id if successful
	CreatePolyGeofence(c PolyGeofence) (string, error)

	// Update an existed poly geofence
	UpdatePolyGeofence(c PolyGeofence) error

	// Alarms
	QueryStatus(monitoredPerson string, fenceIds []string) (QueryStatus, error)
	GetHistoryAlarms(monitoredPerson string, fenceIds []string) (HistoryAlarms, error)
	BatchGetHistoryAlarms(input *BatchGetHistoryAlarmsRequest) (BatchHistoryAlarmsResp, error)
	GetStayPoints(input *GetStayPointsRequest) (StayPoints, error)
	UnmarshalAlarmNotification(content []byte) (*AlarmNotification, error)

	//Entity
	AddEntity(EntityName string, EntityDesc string) error
	UpdateEntity(EntityName string, EntityDesc string) error
	DeleteEntity(EntityName string) error
	ListEntity(collectionId string, CoordTypeOutput string, PageIndex int, pageSize int) (int, ListEntityStruct)
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
