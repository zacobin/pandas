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
package engines

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cloustone/pandas/lbs"
	"github.com/sirupsen/logrus"
)

type baiduLbsRequest struct {
	AK        string `json:"ak,noempty"`
	ServiceId string `json:"service_id,noempty"`
	SN        string `json:"sn,omitempty"`
}

type baiduLbsResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type baiduTrackPoint struct {
	EntityName     string  `json:"entity_name,noempty"`
	Latitude       float64 `json:"latitude,noempty"`
	Longitude      float64 `json:"longitude,noempty"`
	LocTime        int64   `json:"loc_time,noempty"`
	CoordTypeInput string  `json:"coord_type_input,noempty"`
	Speed          float64 `json:"speed"`
	Direction      float64 `json:"direction"`
	Height         float64 `json:"height"`
	Radius         float64 `json:"radius"`
	ObjectName     string  `json:"object_name"`
	ColumnKey      string  `json:"column-key"` // TODO:
}

type baiduAddTrackPointRequest struct {
	baiduLbsRequest
	baiduTrackPoint
}

type baiduAddTrackPointsRequest struct {
	baiduLbsRequest
	Points []baiduTrackPoint `json:"point_list,noempty"`
}

type baiduLbsManager struct {
	accessKey string
	serviceId string
}

func NewBaiduLbsEngine(locationServingOptions *lbs.LocationServingOptions) lbs.LocationEngine {
	return &baiduLbsManager{
		accessKey: locationServingOptions.AK,
		serviceID: locationServingOptions.ServiceId,
	}
}

func (b *baiduLbsManager) AddTrackPoint(point lbs.TrackPoint) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/track/addpoint"
	baiduReq := baiduAddTrackPointRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		baiduTrackPoint: baiduTrackPoint{
			EntityName: point.EntityName,
			Latitude:   point.Latitude,
			Longitude:  point.Longitude,
		},
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	locTime := int(getUnixTimeStamp(point.Time))
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "entity_name": {point.EntityName}, "longitude": {fmt.Sprint(point.Longitude)},
		"latitude": {fmt.Sprint(point.Latitude)}, "loc_time": {fmt.Sprint(locTime)}, "coord_type_input": {"bd09ll"}})
	if err != nil {
		logrus.WithError(err).Errorln("AddTrackPoint failed:", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugf("AddTrackPoint resp:%s", string(body))
	return
}

func (b *baiduLbsManager) AddTrackPoints(points []lbs.TrackPoint) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/track/addpoints"
	baiduReq := baiduAddTrackPointsRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		Points: []baiduTrackPoint{},
	}
	for _, point := range points {
		baiduReq.Points = append(baiduReq.Points, baiduTrackPoint{
			EntityName:     point.EntityName,
			Latitude:       point.Latitude,
			Longitude:      point.Longitude,
			CoordTypeInput: string(point.CoordType),
		})
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	baiduReq.SN = sn

	pointList, _ := json.Marshal(&baiduReq.Points)

	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "point_list": {string(pointList)}})
	if err != nil {
		logrus.WithError(err).Errorln("AddTrackPoint failed:", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugf("AddTrackPoint resp:%s", string(body))
	return
}

type baiduCircleGeofence struct {
	Name             string        `json:"fence_name"`
	MonitoredObjects string        `json:"monitored_persion"`
	Longitude        float64       `json:"longitude,noempty"`
	Latitude         float64       `json:"latitude,noempty"`
	Radius           float64       `json:"radius,noempty"`
	CoordType        lbs.CoordType `json:"coord_type,noempty"`
	Denoise          int           `json:"denoise"`
}

type baiduCreateCircleGeofenceRequest struct {
	baiduLbsRequest
	baiduCircleGeofence
}

type baiduCreateCircleGeofenceResponse struct {
	baiduLbsResponse
	FenceID string `json:"fence_id"`
}

func (b *baiduLbsManager) CreateCircleGeofence(c lbs.CircleGeofence) (string, error) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/createcirclefence"
	id := ""

	baiduReq := baiduCreateCircleGeofenceRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		baiduCircleGeofence: baiduCircleGeofence{
			Name:             c.Name,
			MonitoredObjects: c.MonitoredObjects,
			Longitude:        c.Longitude,
			Latitude:         c.Latitude,
			Radius:           c.Radius,
			CoordType:        c.CoordType,
			Denoise:          c.Denoise,
		},
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	logrus.Debugln("baiduReq:", baiduReq)
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "fence_name": {c.Name}, "monitored_persion": {c.MonitoredObjects},
		"longitude": {fmt.Sprint(c.Longitude)}, "latitude": {fmt.Sprint(c.Latitude)}, "radius": {fmt.Sprint(c.Radius)}, "coord_type": {string(c.CoordType)}, "denoise": {fmt.Sprint(c.Denoise)}})
	if err != nil {
		logrus.WithError(err).Errorln("create circle geofence failed:", err)
		return id, err
	}

	logrus.Debugln("coord_type:", string(c.CoordType))
	rsp := baiduCreateCircleGeofenceResponse{}
	if resp.StatusCode != http.StatusOK {
		return id, fmt.Errorf("exepected status 200, return %d", resp.StatusCode)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		logrus.Debugf("resp:%s", string(body))
		if err := json.Unmarshal(body, &rsp); err != nil {
			return id, err
		}
		id = fmt.Sprint(rsp.FenceID)
	}
	return id, nil
}

type baiduPolyGeofence struct {
	Name             string        `json:"fence_name"`
	MonitoredObjects string        `json:"monitored_persion"`
	Vertexes         string        `json:"vertexes"`
	CoordType        lbs.CoordType `json:"coord_type,noempty"`
	Denoise          int           `json:"denoise"`
}

type baiduCreatePolyGeofenceRequest struct {
	baiduLbsRequest
	baiduPolyGeofence
}

type baiduCreatePolyGeofenceResponse struct {
	baiduLbsResponse
	FenceID string `json:"fence_id"`
}

func (b *baiduLbsManager) CreatePolyGeofence(c lbs.PolyGeofence) (string, error) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/createpolygonfence"
	id := ""

	baiduReq := baiduCreatePolyGeofenceRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		baiduPolyGeofence: baiduPolyGeofence{
			Name:             c.Name,
			MonitoredObjects: c.MonitoredObjects,
			Vertexes:         c.Vertexes,
			CoordType:        c.CoordType,
			Denoise:          c.Denoise,
		},
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	logrus.Debugln("baiduReq:", baiduReq)
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "fence_name": {c.Name}, "monitored_persion": {c.MonitoredObjects},
		"vertexes": {c.Vertexes}, "coord_type": {string(c.CoordType)}, "denoise": {fmt.Sprint(c.Denoise)}})
	if err != nil {
		logrus.WithError(err).Errorln("create poly geofence failed:", err)
		return id, err
	}

	logrus.Debugln("coord_type:", string(c.CoordType))
	rsp := baiduCreatePolyGeofenceResponse{}
	if resp.StatusCode != http.StatusOK {
		return id, fmt.Errorf("exepected status 200, return %d", resp.StatusCode)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		logrus.Debugf("resp:%s", string(body))
		if err := json.Unmarshal(body, &rsp); err != nil {
			return id, err
		}
		id = fmt.Sprint(rsp.FenceID)
	}
	return id, nil
}

func (b *baiduLbsManager) UpdatePolyGeofence(c lbs.PolyGeofence) error {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/updatepolygonfence"
	baiduReq := baiduCreatePolyGeofenceRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		baiduPolyGeofence: baiduPolyGeofence{
			Name:             c.Name,
			MonitoredObjects: c.MonitoredObjects,
			Vertexes:         c.Vertexes,
			CoordType:        c.CoordType,
			Denoise:          c.Denoise,
		},
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "fence_id": {c.FenceID}, "fence_name": {c.Name}, "monitored_persion": {c.MonitoredObjects},
		"vertexes": {c.Vertexes}, "coord_type": {string(c.CoordType)}, "denoise": {fmt.Sprint(c.Denoise)}})
	if err != nil {
		logrus.WithError(err).Errorln("update ciricle geofence failed:", err)
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugf("resp:%s", string(body))

	return nil
}

func (b *baiduLbsManager) UpdateCircleGeofence(c lbs.CircleGeofence) error {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/updatecirclefence"
	baiduReq := baiduCreateCircleGeofenceRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		baiduCircleGeofence: baiduCircleGeofence{
			Name:             c.Name,
			MonitoredObjects: c.MonitoredObjects,
			Longitude:        c.Longitude,
			Latitude:         c.Latitude,
			Radius:           c.Radius,
			CoordType:        c.CoordType,
			Denoise:          c.Denoise,
		},
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "fence_id": {c.FenceID}, "fence_name": {c.Name}, "monitored_persion": {c.MonitoredObjects},
		"longitude": {fmt.Sprint(c.Longitude)}, "latitude": {fmt.Sprint(c.Latitude)}, "radius": {fmt.Sprint(c.Radius)}, "coord_type": {string(c.CoordType)}, "denoise": {fmt.Sprint(c.Denoise)}})
	if err != nil {
		logrus.WithError(err).Errorln("update ciricle geofence failed:", err)
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugf("resp:%s", string(body))

	return nil
}

type baiduDeleteGeofence struct {
	MonitoredObject string `json:"monitored_person"`
	FenceIDs        string `json:"fence_ids"`
}

type baiduDeleteGeofenceRequest struct {
	baiduLbsRequest
	baiduDeleteGeofence
}

type baiduDeleteGeofenceResponse struct {
	baiduLbsResponse
	FenceIDs []int `json:"fence_ids"`
}

func (b *baiduLbsManager) DeleteGeofence(fenceIDs []string, objects []string) ([]string, error) {
	logrus.Debugln("baidulbs deleteGeofence")
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/delete"
	baiduReq := baiduDeleteGeofenceRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		baiduDeleteGeofence: baiduDeleteGeofence{
			MonitoredObject: strings.Join(objects, ","),
			FenceIDs:        strings.Join(fenceIDs, ","),
		},
	}
	sn, aksnErr := caculateAKSN(baiduYYurl, baiduReq)
	if aksnErr != nil {
		logrus.Errorln("caculateAKSN error:", aksnErr)
	}
	logrus.Debugln("fenceIDs:", fenceIDs)
	logrus.Debugln("objects:", objects)
	logrus.Debugln("sn:", sn)

	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "fence_ids": fenceIDs, "monitored_persion": {baiduReq.MonitoredObject}})
	if err != nil {
		logrus.WithError(err).Errorln("update ciricle geofence failed:", err)
		return nil, err
	}

	rsp := baiduDeleteGeofenceResponse{}
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(body))
	if err := json.Unmarshal(body, &rsp); err != nil {
		return nil, err
	}
	fenceIDsResp := []string{}
	for _, val := range rsp.FenceIDs {
		fenceIDsResp = append(fenceIDsResp, strconv.Itoa(val))
	}
	return fenceIDsResp, nil
}

type baiduListGeofenceRequest struct {
	baiduLbsRequest
	FenceName       string `json:"fence_name"`
	MonitoredObject string `json:"monitored_person"`
}
type baiduListGeofenceResponse struct {
	baiduLbsResponse
	Size   int             `json:"int"`
	Fences []*lbs.Geofence `json:"fences"`
}

func (b *baiduLbsManager) ListGeofence(fenceIDs []string, objects []string) ([]*lbs.Geofence, error) {
	logrus.Debugln("ListGeofence")
	url := "http://yingyan.baidu.com/api/v3/fence/list"
	baiduReq := baiduListGeofenceRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		MonitoredObject: strings.Join(objects, ","),
		FenceName:       strings.Join(fenceIDs, ","),
	}
	sn, _ := caculateAKSN(url, baiduReq)
	if len(fenceIDs) > 0 {
		url = fmt.Sprintf("%s?ak=%s&service_id=%s&fence_ids=%s&sn=%s",
			url, b.accessKey, b.serviceId, strings.Join(fenceIDs, ","), sn)
	} else {
		url = fmt.Sprintf("%s?ak=%s&service_id=%s&monitored_persion=%s&sn=%s",
			url, b.accessKey, b.serviceId, strings.Join(objects, ","), sn)
	}
	logrus.Debugln("url:", url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("list geofence failed")
		return nil, err
	}
	rsp := baiduListGeofenceResponse{}

	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	if resp.StatusCode != http.StatusOK {
		logrus.Debugln("status not 200")
		return nil, fmt.Errorf("exepected status 200, return %d", resp.StatusCode)
	} else {
		if err := json.Unmarshal(data, &rsp); err != nil {
			logrus.Debugln("json unmarshal failed:", err)
			return nil, err
		}
	}
	logrus.Debugln("rsp:", rsp)
	return rsp.Fences, nil
}

type baiduAddObjectRequest struct {
	baiduLbsRequest
	FenceID         string `json:"fence_id,noempty"`
	MonitoredObject string `json:"monitored_person,noempty"`
}

func (b *baiduLbsManager) AddMonitoredObject(fenceID string, objects []string) error {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/addmonitoredperson"

	baiduReq := baiduAddObjectRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		MonitoredObject: strings.Join(objects, ","),
		FenceID:         fenceID,
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)

	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "fence_id": {fenceID}, "monitored_person": objects})
	if err != nil {
		logrus.WithError(err).Errorln("update ciricle geofence failed:", err)
		return err
	}

	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	return nil
}

func (b *baiduLbsManager) RemoveMonitoredObject(fenceID string, objects []string) error {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/deletemonitoredperson"

	baiduReq := baiduAddObjectRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		MonitoredObject: strings.Join(objects, ","),
		FenceID:         fenceID,
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "fence_id": {fenceID}, "monitored_person": objects})
	if err != nil {
		logrus.WithError(err).Errorln("update ciricle geofence failed:", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("exepected status 200, return %d", resp.StatusCode)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	return nil
}

type baiduListMonitoredObjectsRequest struct {
	baiduLbsRequest
	FenceID   int `json:"fence_id,noempty"`
	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`
}

type baiduListMonitoredObjectsResponse struct {
	baiduLbsResponse
	Total           int      `json:"total"`
	Size            int      `json:"size"`
	MonitoredPerson []string `json:"monitored_person"`
}

func (b *baiduLbsManager) ListMonitoredObjects(fenceID string, pageIndex int, pageSize int) (int, []string) {
	url := "http://yingyan.baidu.com/api/v3/fence/listmonitoredperson"
	id, _ := strconv.Atoi(fenceID)

	baiduReq := baiduListMonitoredObjectsRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		FenceID:   id,
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}
	sn, _ := caculateAKSN(url, baiduReq)
	url = fmt.Sprintf("%s?ak=%s&service_id=%s&fence_id=%s&page_index=%d&page_size=%d&sn=%s",
		url, b.accessKey, b.serviceId, fenceID, pageIndex, pageSize, sn)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("list monitored objects failed")
		return -1, nil
	}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	rsp := baiduListMonitoredObjectsResponse{}
	if resp.StatusCode != http.StatusOK {
		return -1, nil
	} else {
		if err := json.Unmarshal(data, &rsp); err != nil {
			return -1, nil
		}
	}
	return rsp.Total, rsp.MonitoredPerson
}

func (b *baiduLbsManager) QueryStatus(monitoredPerson string, fenceIDs []string) (lbs.QueryStatus, error) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/querystatus"

	baiduYYurl = fmt.Sprintf("%s?ak=%s&service_id=%s&monitored_person=%s",
		baiduYYurl, b.accessKey, b.serviceId, monitoredPerson) + "&fence_ids=" + strings.Join(fenceIDs, ",")
	logrus.Debugln("baiduYYurl:", baiduYYurl)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", baiduYYurl, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("list monitored objects failed")
		return BaiduQueryStatusResponse{}, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	rsp := BaiduQueryStatusResponse{}
	if resp.StatusCode != http.StatusOK {
		return BaiduQueryStatusResponse{}, err
	} else {
		if err := json.Unmarshal(data, &rsp); err != nil {
			logrus.Errorln("QueryStatus json unmarshal failed:", err)
			return BaiduQueryStatusResponse{}, err
		}
	}
	logrus.Debugln("rsp:", rsp)

	return rsp, nil
}

func (b *baiduLbsManager) GetHistoryAlarms(monitoredPerson string, fenceIDs []string) (lbs.HistoryAlarms, error) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/historyalarm"

	baiduYYurl = fmt.Sprintf("%s?ak=%s&service_id=%s&monitored_person=%s",
		baiduYYurl, b.accessKey, b.serviceId, monitoredPerson) + "&fence_ids=" + strings.Join(fenceIDs, ",")
	logrus.Debugln("baiduYYurl:", baiduYYurl)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", baiduYYurl, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("list monitored objects failed")
		return BaiduGetHistoryAlarmsResponse{}, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	rsp := BaiduGetHistoryAlarmsResponse{}
	if resp.StatusCode != http.StatusOK {
		return BaiduGetHistoryAlarmsResponse{}, err
	} else {
		if err := json.Unmarshal(data, &rsp); err != nil {
			logrus.Errorln("QueryStatus json unmarshal failed:", err)
			return BaiduGetHistoryAlarmsResponse{}, err
		}
	}
	logrus.Debugln("rsp:", rsp)

	return rsp, nil
}

func (b *baiduLbsManager) BatchGetHistoryAlarms(input *lbs.HistoryAlarmsRequest) (lbs.HistoryAlarms, error) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/fence/batchhistoryalarm"

	startTime := int(getUnixTimeStamp(input.StartTime))
	endTime := int(getUnixTimeStamp(input.EndTime))

	baiduYYurl = fmt.Sprintf("%s?ak=%s&service_id=%s&start_time=%d&end_time=%d&coord_type_output=%s&page_index=%d&page_size=%d",
		baiduYYurl, b.accessKey, b.serviceId, startTime, endTime, input.CoordTypeOutput, input.PageIndex, input.PageSize)
	logrus.Debugln("baiduYYurl:", baiduYYurl)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", baiduYYurl, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("list monitored objects failed")
		return BaiduBatchHistoryAlarmsResp{}, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	rsp := BaiduBatchHistoryAlarmsResp{}
	if resp.StatusCode != http.StatusOK {
		return BaiduBatchHistoryAlarmsResp{}, err
	} else {
		if err := json.Unmarshal(data, &rsp); err != nil {
			logrus.Errorln("QueryStatus json unmarshal failed:", err)
			return BaiduBatchHistoryAlarmsResp{}, err
		}
	}
	logrus.Debugln("rsp:", rsp)

	return rsp, nil
}

func getUnixTimeStamp(strTime string) int64 {
	timeStamp, err := time.Parse("2006-01-02 15:04:05", strTime)
	if err != nil {
		logrus.Errorln("err:", err)
		return 0
	}
	timeStamp = timeStamp.Add(-8 * time.Hour)
	logrus.Debugln("strTIme:", strTime, "now:", time.Now())
	unixTime := timeStamp.Unix()
	logrus.Debugln("unixTime:", unixTime)
	return unixTime
}

type BaiduGetStayPointResp struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Size       int         `json:"size"`
	Total      int         `json:"total"`
	StartPoint lbs.Point   `json:"start_point"`
	EndPoint   lbs.Point   `json:"end_point"`
	Points     []lbs.Point `json:"points"`
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

func (b *baiduLbsManager) GetStayPoints(input *GetStayPointsRequest) (BaiduGetStayPointResp, error) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/track/gettrack"

	startTime := int(getUnixTimeStamp(input.StartTime))
	endTime := int(getUnixTimeStamp(input.EndTime))

	baiduYYurl = fmt.Sprintf("%s?ak=%s&service_id=%s&start_time=%d&end_time=%d&coord_type_output=%s&page_index=%d&page_size=%d&entity_name=%s",
		baiduYYurl, b.accessKey, b.serviceId, startTime, endTime, input.CoordTypeOutput, input.PageIndex, input.PageSize, input.EntityName)

	logrus.Debugln("baiduYYurl:", baiduYYurl)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", baiduYYurl, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("list monitored objects failed")
		return BaiduGetStayPointResp{}, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	rsp := BaiduGetStayPointResp{}
	if resp.StatusCode != http.StatusOK {
		return BaiduGetStayPointResp{}, err
	} else {
		if err := json.Unmarshal(data, &rsp); err != nil {
			logrus.Errorln("QueryStatus json unmarshal failed:", err)
			return BaiduGetStayPointResp{}, err
		}
	}
	logrus.Debugln("rsp:", rsp)

	return rsp, nil
}

type baiduAddEntityRequest struct {
	baiduLbsRequest
	EntityName string `json:"entity_name"`
	EntityDesc string `json:"entity_desc"`
}

type baiduAddEntityResponse struct {
	baiduLbsResponse
}

func (b *baiduLbsManager) AddEntity(entityName string, entityDesc string) error {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/entity/add"

	baiduReq := baiduAddEntityRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		EntityName: entityName,
		EntityDesc: entityDesc,
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "entity_name": {entityName}, "entity_desc": {entityDesc}})
	if err != nil {
		logrus.WithError(err).Errorln("update ciricle geofence failed:", err)
		return err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("exepected status 200, return %d", resp.StatusCode)
	}
	return nil
}

type baiduListEntityRequest struct {
	baiduLbsRequest
	CoordTypeOutput string `json:"coord_type_output"`
	PageIndex       int  `json:"page_index"`
	PageSize        int  `json:"page_size"`
}

type baiduListEntityResponse struct {
}

func (b *baiduLbsManager) ListEntity(collectionID string, CoordTypeOutput string, PageIndex int, pageSize int) (int, lbs.ListEntityStruct) {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/entity/list"

	baiduReq := baiduListEntityRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		CoordTypeOutput: CoordTypeOutput,
		PageIndex:       PageIndex,
		PageSize:        pageSize,
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	baiduYYurl = fmt.Sprintf("%s?ak=%s&service_id=%s&coord_type_output=%s&page_index=%d&page_size=%d&sn=%s",
		baiduYYurl, b.accessKey, b.serviceId, CoordTypeOutput, PageIndex, pageSize, sn)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", baiduYYurl, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("list monitored objects failed")
		return -1, baiduListEntityStruct{}
	}
	rsp := baiduListEntityStruct{}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))

	if resp.StatusCode != http.StatusOK {
		return -1, baiduListEntityStruct{}
	} else {
		if err := json.Unmarshal(data, &rsp); err != nil {
			logrus.Errorln("list entity unmarshal failed:", err)
			return -1, baiduListEntityStruct{}
		}
		logrus.Debugln("total:", rsp.Total)
		logrus.Debugln("rsp:", rsp)
	}

	return rsp.Total, rsp
}

type baiduUpdateEntityRequest struct {
	baiduLbsRequest
	EntityName string `json:"entity_name"`
	EntityDesc string `json:"entity_desc"`
}

type baiduUpdateEntityResponse struct {
	baiduLbsResponse
}

func (b *baiduLbsManager) UpdateEntity(entityName string, entityDesc string) error {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/entity/update"

	baiduReq := baiduUpdateEntityRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		EntityName: entityName,
		EntityDesc: entityDesc,
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "entity_name": {entityName}, "entity_desc": {entityDesc}})
	if err != nil {
		logrus.WithError(err).Errorln("UpdateEntity failed:", err)
		return err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("exepected status 200, return %d", resp.StatusCode)
	}
	return nil
}

type baiduDeleteEntityRequest struct {
	baiduLbsRequest
	EntityName string `json:"entity_name"`
}

type baiduDeleteEntityResponse struct {
	baiduLbsResponse
}

func (b *baiduLbsManager) DeleteEntity(entityName string) error {
	baiduYYurl := "http://yingyan.baidu.com/api/v3/entity/delete"

	baiduReq := baiduDeleteEntityRequest{
		baiduLbsRequest: baiduLbsRequest{
			AK:        b.accessKey,
			ServiceId: b.serviceId,
		},
		EntityName: entityName,
	}
	sn, _ := caculateAKSN(baiduYYurl, baiduReq)
	resp, err := http.PostForm(baiduYYurl, url.Values{"ak": {b.accessKey}, "service_id": {b.serviceId},
		"sn": {sn}, "entity_name": {entityName}})
	if err != nil {
		logrus.WithError(err).Errorln("DeleteEntity failed:", err)
		return err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("resp:", string(data))
	if resp.StatusCode != http.StatusOK {
		logrus.Errorln("exepected status 200, return ", resp.StatusCode)
		return fmt.Errorf("exepected status 200, return %d", resp.StatusCode)
	}
	return nil
}

type baiduLocationPoint struct {
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Radius     int     `json:"radius"`
	CoordType  string  `json:"coord_type"`
	LocTime    string  `json:"loc_time"`
	CreateTime string  `json:"create_time"`
}

type baiduAlarm struct {
	FenceID          string             `json:"fence_id,noempty"`
	FenceName        string             `json:"fence_name,noempty"`
	MonitoredObjects string             `json:"monitored_person"`
	Action           string             `json:"action"`
	AlarmPoint       baiduLocationPoint `json:"alarm_point"`
	PrePoint         baiduLocationPoint `json:"pre_point"`
}

type baiduAlarmNotification struct {
	Type      int           `json:"type"`
	ServiceId int           `json:"service_id"`
	Alarms    []*baiduAlarm `json:"content"`
}

func (b *baiduLbsManager) UnmarshalAlarmNotification(content []byte) (*lbs.AlarmNotification, error) {
	logrus.Debugf("unmarshal baidu alarm notification")

	n := baiduAlarmNotification{}
	if err := json.Unmarshal(content, &n); err != nil {
		logrus.WithError(err).Errorf("unmarshal baidu alarm notifcation failed")
		return nil, err
	}
	alarmNotify := &AlarmNotification{
		Type:      n.Type,
		ServiceId: strconv.Itoa(n.ServiceId),
		Alarms:    []*Alarm{},
	}
	for _, alarm := range n.Alarms {
		alarmNotify.Alarms = append(alarmNotify.Alarms, &Alarm{
			FenceID:          alarm.FenceID,
			FenceName:        alarm.FenceName,
			MonitoredObjects: strings.Split(alarm.MonitoredObjects, ","),
			Action:           alarm.Action,
			AlarmPoint: AlarmPoint{
				Longitude:  alarm.AlarmPoint.Longitude,
				Latitude:   alarm.AlarmPoint.Latitude,
				Radius:     alarm.AlarmPoint.Radius,
				CoordType:  alarm.AlarmPoint.CoordType,
				LocTime:    alarm.AlarmPoint.LocTime,
				CreateTime: alarm.AlarmPoint.CreateTime,
			},
			PrePoint: AlarmPoint{
				Longitude:  alarm.PrePoint.Longitude,
				Latitude:   alarm.PrePoint.Latitude,
				Radius:     alarm.PrePoint.Radius,
				CoordType:  alarm.PrePoint.CoordType,
				LocTime:    alarm.PrePoint.LocTime,
				CreateTime: alarm.PrePoint.CreateTime,
			},
		})
	}
	return alarmNotify, nil
}
