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
package proxy

import (
	"fmt"

	logr "github.com/sirupsen/logrus"
)

type Proxy struct {
	engineName string
	repo       Repository
	engine     Engine
}

func NewProxy(locationServingOptions *LocationServingOptions) *Proxy {
	return &Proxy{
		engine:     newBaiduLbsEngine(locationServingOptions),
		engineName: locationServingOptions.Provider,
		repo:       NewRepository(),
	}
}

func (p *Proxy) AddTrackPoint(point TrackPoint) {
	p.engine.AddTrackPoint(point)
}

func (p *Proxy) AddTrackPoints(points []TrackPoint) {
	p.engine.AddTrackPoints(points)
}

func (p *Proxy) CreateCircleGeofence(c CircleGeofence) (string, error) {
	return p.engine.CreateCircleGeofence(c)
}

func (p *Proxy) UpdateCircleGeofence(c CircleGeofence) error {
	return p.engine.UpdateCircleGeofence(c)
}

func (p *Proxy) DeleteGeofence(fenceIds []string, objects []string) ([]string, error) {
	return p.engine.DeleteGeofence(fenceIds, objects)
}

func (p *Proxy) ListGeofence(fenceIds []string, objects []string) ([]*Geofence, error) {
	return p.engine.ListGeofence(fenceIds, objects)
}

func (p *Proxy) AddMonitoredObject(fenceId string, objects []string) error {
	return p.engine.AddMonitoredObject(fenceId, objects)
}

func (p *Proxy) RemoveMonitoredObject(fenceId string, objects []string) error {
	return p.engine.RemoveMonitoredObject(fenceId, objects)
}

func (p *Proxy) ListMonitoredObjects(fenceId string, pageIndex int, pageSize int) (int, []string) {
	return p.engine.ListMonitoredObjects(fenceId, pageIndex, pageSize)
}

func (p *Proxy) CreatePolyGeofence(c PolyGeofence) (string, error) {
	return p.engine.CreatePolyGeofence(c)
}

func (p *Proxy) UpdatePolyGeofence(c PolyGeofence) error {
	return p.engine.UpdatePolyGeofence(c)
}

// Alarms
func getMonitoredStatus(fenceStatus BaiduQueryStatusResponse) *QueryStatus {
	rsp := &QueryStatus{
		Status:  int32(fenceStatus.Status),
		Message: fenceStatus.Message,
		Size:    int32(fenceStatus.Size),
	}
	for _, mpVal := range fenceStatus.MonitoredStatuses {
		monitoredStatus := MonitoredStatus{
			FenceId:         mpVal.FenceId,
			MonitoredStatus: mpVal.MonitoredStatus,
		}
		rsp.MonitoredStatuses = append(rsp.MonitoredStatuses, monitoredStatus)
	}
	logr.Debugln("get monitered status:", rsp)
	return rsp
}

func (p *Proxy) QueryStatus(monitoredPerson string, fenceIds []string) (*QueryStatus, error) {

	status, err := p.engine.QueryStatus(monitoredPerson, fenceIds)
	rsp := getMonitoredStatus(status)
	return rsp, err
}

func (p *Proxy) GetHistoryAlarms(monitoredPerson string, fenceIds []string) (*HistoryAlarms, error) {
	alarms, err := p.engine.GetHistoryAlarms(monitoredPerson, fenceIds)
	rsp := getHistoryAlarmPoint(alarms)
	return rsp, err
}
func getHistoryAlarmPoint(alarmPoint BaiduGetHistoryAlarmsResponse) *HistoryAlarms {
	rsp := &HistoryAlarms{
		Status:  int32(alarmPoint.Status),
		Message: alarmPoint.Message,
		Size:    int32(alarmPoint.Size),
	}
	for _, haVal := range alarmPoint.Alarms {

		alarm := &Alarm{
			FenceId:          haVal.FenceId,
			FenceName:        haVal.FenceName,
			MonitoredObjects: haVal.MonitoredPerson,
			Action:           haVal.Action,
			AlarmPoint: AlarmPoint{
				Longitude:  haVal.AlarmPoint.Longitude,
				Latitude:   haVal.AlarmPoint.Latitude,
				Radius:     haVal.AlarmPoint.Radius,
				CoordType:  haVal.AlarmPoint.CoordType,
				LocTime:    haVal.AlarmPoint.LocTime,
				CreateTime: haVal.AlarmPoint.CreateTime,
			},
			PrePoint: AlarmPoint{
				Longitude:  haVal.AlarmPoint.Longitude,
				Latitude:   haVal.AlarmPoint.Latitude,
				Radius:     haVal.AlarmPoint.Radius,
				CoordType:  haVal.AlarmPoint.CoordType,
				LocTime:    haVal.AlarmPoint.LocTime,
				CreateTime: haVal.AlarmPoint.CreateTime,
			},
		}

		rsp.Alarms = append(rsp.Alarms, alarm)
	}
	logr.Debugln("getHistoryAlarmPoint:", rsp)
	return rsp
}

func (p *Proxy) BatchGetHistoryAlarms(input *BatchGetHistoryAlarmsRequest) (*HistoryAlarms, error) {
	historyAlarms, err := p.engine.BatchGetHistoryAlarms(input)
	rsp := getBatchHistoryAlarmPoint(historyAlarms)
	return rsp, err
}
func getBatchHistoryAlarmPoint(historyAlarms BaiduBatchHistoryAlarmsResp) *HistoryAlarms {
	rsp := &HistoryAlarms{
		Status:  int32(historyAlarms.Status),
		Message: historyAlarms.Message,
		Size:    int32(historyAlarms.Size),
		Total:   int32(historyAlarms.Total),
	}
	for _, haVal := range historyAlarms.Alarms {

		alarm := &Alarm{
			FenceId:          haVal.FenceId,
			FenceName:        haVal.FenceName,
			MonitoredObjects: haVal.MonitoredPerson,
			Action:           haVal.Action,
			AlarmPoint: AlarmPoint{
				Longitude:  haVal.AlarmPoint.Longitude,
				Latitude:   haVal.AlarmPoint.Latitude,
				Radius:     haVal.AlarmPoint.Radius,
				CoordType:  haVal.AlarmPoint.CoordType,
				LocTime:    haVal.AlarmPoint.LocTime,
				CreateTime: haVal.AlarmPoint.CreateTime,
			},
			PrePoint: AlarmPoint{
				Longitude:  haVal.AlarmPoint.Longitude,
				Latitude:   haVal.AlarmPoint.Latitude,
				Radius:     haVal.AlarmPoint.Radius,
				CoordType:  haVal.AlarmPoint.CoordType,
				LocTime:    haVal.AlarmPoint.LocTime,
				CreateTime: haVal.AlarmPoint.CreateTime,
			},
		}

		rsp.Alarms = append(rsp.Alarms, alarm)
	}
	logr.Debugln("getHistoryAlarmPoint:", rsp)
	return rsp
}

func (p *Proxy) GetStayPoints(input *GetStayPointsRequest) (*StayPoints, error) {
	stay, err := p.engine.GetStayPoints(input)
	rsp := getGrpcStayPoints(stay)
	return rsp, err

}
func getGrpcStayPoints(stayPoints BaiduGetStayPointResp) *StayPoints {
	rsp := &StayPoints{
		Status:  int32(stayPoints.Status),
		Message: stayPoints.Message,
		Size:    int32(stayPoints.Size),
		Total:   int32(stayPoints.Total),
		StartPoint: &Point{
			Latitude:  stayPoints.StartPoint.Latitude,
			Longitude: stayPoints.StartPoint.Longitude,
			CoordType: stayPoints.StartPoint.CoordType,
			LocTime:   fmt.Sprint(stayPoints.StartPoint.LocTime),
		},
		EndPoint: &Point{
			Latitude:  stayPoints.EndPoint.Latitude,
			Longitude: stayPoints.EndPoint.Longitude,
			CoordType: stayPoints.EndPoint.CoordType,
			LocTime:   fmt.Sprint(stayPoints.EndPoint.LocTime),
		},
	}
	for _, val := range stayPoints.Points {
		point := &Point{
			Latitude:  val.Latitude,
			Longitude: val.Longitude,
			CoordType: val.CoordType,
			LocTime:   fmt.Sprint(val.LocTime),
		}
		rsp.Points = append(rsp.Points, point)
	}
	logr.Debugln("getGrpcStayPoints rsp:", rsp)
	return rsp
}

func (p *Proxy) UnmarshalAlarmNotification(content []byte) (*AlarmNotification, error) {
	return p.engine.UnmarshalAlarmNotification(content)
}

//Entity
func (p *Proxy) AddEntity(entityName string, entityDesc string) error {
	return p.engine.AddEntity(entityName, entityDesc)
}

func (p *Proxy) UpdateEntity(entityName string, entityDesc string) error {
	return p.engine.UpdateEntity(entityName, entityDesc)
}

func (p *Proxy) DeleteEntity(entityName string) error {
	return p.engine.DeleteEntity(entityName)
}

func (p *Proxy) ListEntity(collectionId string, coordTypeOutput string, pageIndex int32, pageSize int32) (int, baiduListEntityStruct) {
	return p.engine.ListEntity(collectionId, coordTypeOutput, pageIndex, pageSize)
}
