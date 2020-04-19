// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"

	"github.com/cloustone/pandas/lbs"
	"github.com/go-kit/kit/endpoint"
)

func listCollectionsEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listCollectionsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		saved, err := svc.ListCollections(ctx, req.token)
		if err != nil {
			return nil, err
		}

		res := listCollectionsRes{}
		for _, product := range saved {
			res.Products = append(res.Products, product)
		}
		return res, nil

	}
}

func createCircleGeofenceEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createCircleGeofenceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		fence := &lbs.CircleGeofence{
			Name:      req.fence.Name,
			Longitude: req.fence.Longitude,
			Latitude:  req.fence.Latitude,
			Radius:    req.fence.Radius,
			CoordType: req.fence.CoordType,
			Denoise:   req.fence.Denoise,
			FenceID:   req.fence.FenceID,
		}
		for _, object := range req.fence.MonitoredObjects {
			fence.MonitoredObjects = append(fence.MonitoredObjects, object)
		}

		saved, err := svc.CreateCircleGeofence(ctx, req.token, req.projectID, fence)
		if err != nil {
			return nil, err
		}

		res := createCircleGeofenceRes{
			fenceID: saved,
		}
		return res, nil
	}
}

func updateCircleGeofenceEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateCircleGeofenceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		fence := &lbs.CircleGeofence{
			Name:      req.fence.Name,
			Longitude: req.fence.Longitude,
			Latitude:  req.fence.Latitude,
			Radius:    req.fence.Radius,
			CoordType: req.fence.CoordType,
			Denoise:   req.fence.Denoise,
			FenceID:   req.fence.FenceID,
		}
		for _, object := range req.fence.MonitoredObjects {
			fence.MonitoredObjects = append(fence.MonitoredObjects, object)
		}
		err := svc.UpdateCircleGeofence(ctx, req.token, req.projectID, fence)
		if err != nil {
			res := updateCircleGeofenceRes{
				updated: false,
			}
			return res, err
		}

		res := updateCircleGeofenceRes{
			updated: true,
		}
		return res, nil
	}
}

func deleteGeofenceEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteGeofenceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		err := svc.DeleteGeofence(ctx, req.token, req.projectID, req.fenceIDs, req.objects)
		if err != nil {
			res := deleteGeofenceRes{
				deleted: false,
			}
			return res, err
		}

		res := deleteGeofenceRes{
			deleted: true,
		}
		return res, nil
	}
}

func listGeofencesEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listGeofencesReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		saved, err := svc.ListGeofences(ctx, req.token, req.projectID, req.fenceIDs, req.objects)
		if err != nil {
			return nil, err
		}
		res := listGeofencesRes{}

		for _, f := range saved {
			fence := &Geofence{
				FenceID:   f.FenceID,
				FenceName: f.FenceName,
				//MonitoredObject: f.MonitoredObject,
				Shape:      f.Shape,
				Longitude:  f.Longitude,
				Latitude:   f.Latitude,
				Radius:     f.Radius,
				CoordType:  f.CoordType,
				Denoise:    f.Denoise,
				CreateTime: f.CreateTime,
				UpdateTime: f.UpdateTime,
			}
			for _, vtx := range f.Vertexes {
				vertexe := &Vertexe{
					Latitude:  vtx.Latitude,
					Longitude: vtx.Longitude,
				}
				fence.Vertexes = append(fence.Vertexes, vertexe)
			}
			res.fenceList = append(res.fenceList, fence)
		}
		return res, nil
	}
}

func addMonitoredObjectEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addMonitoredObjectReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		err := svc.AddMonitoredObject(ctx, req.token, req.projectID, req.fenceID, req.objects)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func removeMonitoredObjectEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(removeMonitoredObjectReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		err := svc.RemoveMonitoredObject(ctx, req.token, req.projectID, req.fenceID, req.objects)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func listMonitoredObjectsEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listMonitoredObjectsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		t, o, err := svc.ListMonitoredObjects(ctx, req.token, req.projectID, req.fenceID, req.pageIndex, req.pageSize)
		if err != nil {
			return nil, err
		}

		res := listMonitoredObjectsRes{
			total:   t,
			objects: o,
		}
		for _, object := range o {
			res.objects = append(res.objects, object)
		}
		return res, nil
	}
}

func createPolyGeofenceEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createPolyGeofenceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		fence := &lbs.PolyGeofence{
			Name:      req.fence.Name,
			Vertexes:  req.fence.Vertexes,
			CoordType: req.fence.CoordType,
			Denoise:   req.fence.Denoise,
			FenceID:   req.fence.FenceID,
		}
		for _, object := range req.fence.MonitoredObjects {
			fence.MonitoredObjects = append(fence.MonitoredObjects, object)
		}
		saved, err := svc.CreatePolyGeofence(ctx, req.token, req.projectID, fence)
		if err != nil {
			return nil, err
		}

		res := createPolyGeofenceRes{
			fenceID: saved,
		}
		return res, nil
	}
}

func updatePolyGeofenceEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updatePolyGeofenceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		fence := &lbs.PolyGeofence{
			Name:      req.fence.Name,
			Vertexes:  req.fence.Vertexes,
			CoordType: req.fence.CoordType,
			Denoise:   req.fence.Denoise,
			FenceID:   req.fence.FenceID,
		}
		for _, object := range req.fence.MonitoredObjects {
			fence.MonitoredObjects = append(fence.MonitoredObjects, object)
		}
		err := svc.UpdatePolyGeofence(ctx, req.token, req.projectID, fence)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func getFenceIDsEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getFenceIDsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		saved, err := svc.GetFenceIDs(ctx, req.token, req.projectID)
		if err != nil {
			return nil, err
		}
		res := getFenceIDsRes{
			fenceIDs: saved,
		}
		return res, nil
	}
}

func queryStatusEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(queryStatusReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		saved, err := svc.QueryStatus(ctx, req.token, req.projectID, req.monitoredPerson, req.fenceIDs)
		if err != nil {
			return nil, err
		}
		res := queryStatusRes{
			Status:  saved.Status,
			Message: saved.Message,
			Size:    saved.Size,
		}
		for _, m := range saved.MonitoredStatuses {
			status := MonitoredStatus{
				FenceID:         m.FenceID,
				MonitoredStatus: m.MonitoredStatus,
			}
			res.MonitoredStatuses = append(res.MonitoredStatuses, status)
		}
		return res, nil
	}
}

func getHistoryAlarmsEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getHistoryAlarmsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		saved, err := svc.GetHistoryAlarms(ctx, req.token, req.projectID, req.monitoredPerson, req.fenceIDs)
		if err != nil {
			return nil, err
		}
		res := getHistoryAlarmsRes{
			Status:  saved.Status,
			Message: saved.Message,
			Total:   saved.Total,
			Size:    saved.Size,
		}
		for _, a := range saved.Alarms {
			alarm := Alarm{
				FenceID:   a.FenceID,
				FenceName: a.FenceName,
				Action:    a.Action,
				AlarmPoint: AlarmPoint{
					Longitude:  a.AlarmPoint.Longitude,
					Latitude:   a.AlarmPoint.Latitude,
					Radius:     a.AlarmPoint.Radius,
					CoordType:  a.AlarmPoint.CoordType,
					LocTime:    a.AlarmPoint.LocTime,
					CreateTime: a.AlarmPoint.CreateTime,
				},
				PrePoint: AlarmPoint{
					Longitude:  a.PrePoint.Longitude,
					Latitude:   a.PrePoint.Latitude,
					Radius:     a.PrePoint.Radius,
					CoordType:  a.PrePoint.CoordType,
					LocTime:    a.PrePoint.LocTime,
					CreateTime: a.PrePoint.CreateTime,
				},
			}
			for _, m := range a.MonitoredObjects {
				alarm.MonitoredObjects = append(alarm.MonitoredObjects, m)
			}
			res.Alarms = append(res.Alarms, &alarm)
		}
		return res, nil
	}
}

func batchGetHistoryAlarmsEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(batchGetHistoryAlarmsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		tmp := &lbs.BatchGetHistoryAlarmsRequest{
			EndTime:         req.input.EndTime,
			StartTime:       req.input.StartTime,
			PageIndex:       req.input.PageIndex,
			PageSize:        req.input.PageSize,
			CoordTypeOutput: req.input.CoordTypeOutput,
		}
		saved, err := svc.BatchGetHistoryAlarms(ctx, req.token, req.projectID, tmp)
		if err != nil {
			return nil, err
		}
		res := batchGetHistoryAlarmsRes{
			Status:  saved.Status,
			Message: saved.Message,
			Total:   saved.Total,
			Size:    saved.Size,
		}
		for _, a := range saved.Alarms {
			alarm := Alarm{
				FenceID:   a.FenceID,
				FenceName: a.FenceName,
				Action:    a.Action,
				AlarmPoint: AlarmPoint{
					Longitude:  a.AlarmPoint.Longitude,
					Latitude:   a.AlarmPoint.Latitude,
					Radius:     a.AlarmPoint.Radius,
					CoordType:  a.AlarmPoint.CoordType,
					LocTime:    a.AlarmPoint.LocTime,
					CreateTime: a.AlarmPoint.CreateTime,
				},
				PrePoint: AlarmPoint{
					Longitude:  a.PrePoint.Longitude,
					Latitude:   a.PrePoint.Latitude,
					Radius:     a.PrePoint.Radius,
					CoordType:  a.PrePoint.CoordType,
					LocTime:    a.PrePoint.LocTime,
					CreateTime: a.PrePoint.CreateTime,
				},
			}
			for _, m := range a.MonitoredObjects {
				alarm.MonitoredObjects = append(alarm.MonitoredObjects, m)
			}
			res.Alarms = append(res.Alarms, &alarm)
		}
		return res, nil
	}
}

func getStayPointsEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getStayPointsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		tmp := &lbs.GetStayPointsRequest{
			EndTime:         req.input.EndTime,
			StartTime:       req.input.StartTime,
			PageIndex:       req.input.PageIndex,
			PageSize:        req.input.PageSize,
			CoordTypeOutput: req.input.CoordTypeOutput,
			EntityName:      req.input.EntityName,
		}
		for _, id := range req.input.FenceIDs {
			tmp.FenceIDs = append(tmp.FenceIDs, id)
		}
		saved, err := svc.GetStayPoints(ctx, req.token, req.projectID, tmp)
		if err != nil {
			return nil, err
		}
		res := getStayPointsRes{
			Status:   saved.Status,
			Message:  saved.Message,
			Total:    saved.Total,
			Size:     saved.Size,
			Distance: saved.Distance,
			EndPoint: &Point{
				Longitude: saved.EndPoint.Latitude,
				Latitude:  saved.EndPoint.Latitude,
				CoordType: saved.EndPoint.CoordType,
				LocTime:   saved.EndPoint.LocTime,
			},
			StartPoint: &Point{
				Longitude: saved.StartPoint.Latitude,
				Latitude:  saved.StartPoint.Latitude,
				CoordType: saved.StartPoint.CoordType,
				LocTime:   saved.StartPoint.LocTime,
			},
		}
		for _, point := range saved.Points {
			p := &Point{
				Longitude: point.Latitude,
				Latitude:  point.Latitude,
				CoordType: point.CoordType,
				LocTime:   point.LocTime,
			}

			res.Points = append(res.Points, p)
		}
		return res, nil
	}
}

func notifyAlarmsEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(notifyAlarmsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.NotifyAlarms(ctx, req.token, req.projectID, req.content)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func getFenceUserIDEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getFenceUserIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		saved, err := svc.GetFenceUserID(ctx, req.token, req.fenceID)
		if err != nil {
			return nil, err
		}
		res := getFenceUserIDRes{
			UserID: saved,
		}
		return res, nil
	}
}

func addEntityEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addEntityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.AddEntity(ctx, req.token, req.projectID, req.entityName, req.entityDesc)
		if err != nil {
			return nil, err
		}
		res := addEntityRes{
			Successed: true,
		}
		return res, nil
	}
}

func deleteEntityEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteEntityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.DeleteEntity(ctx, req.token, req.projectID, req.entityName)
		if err != nil {
			return nil, err
		}
		res := deleteEntityRes{
			Successed: true,
		}
		return res, nil
	}
}

func updateEntityEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateEntityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.UpdateEntity(ctx, req.token, req.projectID, req.entityName, req.entityDesc)
		if err != nil {
			return nil, err
		}
		res := updateEntityRes{
			Successed: true,
		}
		return res, nil
	}
}

func listEntityEndpoint(svc lbs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listEntityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		total, saved, err := svc.ListEntity(ctx, req.token, req.projectID, req.coordTypeOutput, req.pageIndex, req.pageSize)
		if err != nil {
			return nil, err
		}
		res := listEntityRes{
			Total: total,
		}
		for _, s := range saved {
			info := &EntityInfo{
				EntityName: s.EntityName,
				Latitude:   s.Latitude,
				Longitude:  s.Longitude,
			}
			res.EntityInfos = append(res.EntityInfos, info)
		}
		return res, nil
	}
}
