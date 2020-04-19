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
package providers

/*
// Alarms
func getMonitoredStatus(fenceStatus lbs.QueryStatusResponse) *QueryStatus {
	rsp := &QueryStatus{
		Status:  int(fenceStatus.Status),
		Message: fenceStatus.Message,
		Size:    int(fenceStatus.Size),
	}
	for _, mpVal := range fenceStatus.MonitoredStatuses {
		monitoredStatus := MonitoredStatus{
			FenceID:         mpVal.FenceID,
			MonitoredStatus: mpVal.MonitoredStatus,
		}
		rsp.MonitoredStatuses = append(rsp.MonitoredStatuses, monitoredStatus)
	}
	logr.Debugln("get monitered status:", rsp)
	return rsp
}

func (p *Proxy) QueryStatus(monitoredPerson string, fenceIDs []string) (*QueryStatus, error) {

	status, err := p.engine.QueryStatus(monitoredPerson, fenceIDs)
	rsp := getMonitoredStatus(status)
	return rsp, err
}

func (p *Proxy) GetHistoryAlarms(monitoredPerson string, fenceIDs []string) (*HistoryAlarms, error) {
	alarms, err := p.engine.GetHistoryAlarms(monitoredPerson, fenceIDs)
	rsp := getHistoryAlarmPoint(alarms)
	return rsp, err
}
func getHistoryAlarmPoint(alarmPoint BaiduGetHistoryAlarmsResponse) *HistoryAlarms {
	rsp := &HistoryAlarms{
		Status:  int(alarmPoint.Status),
		Message: alarmPoint.Message,
		Size:    int(alarmPoint.Size),
	}
	for _, haVal := range alarmPoint.Alarms {

		alarm := &Alarm{
			FenceID:          haVal.FenceID,
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
		Status:  int(historyAlarms.Status),
		Message: historyAlarms.Message,
		Size:    int(historyAlarms.Size),
		Total:   int(historyAlarms.Total),
	}
	for _, haVal := range historyAlarms.Alarms {

		alarm := &Alarm{
			FenceID:          haVal.FenceID,
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
		Status:  int(stayPoints.Status),
		Message: stayPoints.Message,
		Size:    int(stayPoints.Size),
		Total:   int(stayPoints.Total),
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

func (p *Proxy) HandleAlarmNotification(content []byte) (*AlarmNotification, error) {
	return p.engine.HandleAlarmNotification(content)
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

func (p *Proxy) ListEntity(collectionID string, coordTypeOutput string, pageIndex int, pageSize int) (int, baiduListEntityStruct) {
	return p.engine.ListEntity(collectionID, coordTypeOutput, pageIndex, pageSize)
}
*/
