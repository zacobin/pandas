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

// Package api contains implementation of lbs service HTTP API.

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/cloustone/pandas/lbs"
	log "github.com/cloustone/pandas/pkg/logger"
)

var _ lbs.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    lbs.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc lbs.Service, logger log.Logger) lbs.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) ListCollections(ctx context.Context, token string) (products []string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_collections for token %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListCollections(ctx, token)
}

// Geofence
func (lm *loggingMiddleware) CreateCircleGeofence(ctx context.Context, token string, projectID string, fence *lbs.CircleGeofence) (fenceID string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_circle_geofence for token %s and project %s and fanceName %s took %s to complete", token, projectID, fence.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreateCircleGeofence(ctx, token, projectID, fence)
}

func (lm *loggingMiddleware) UpdateCircleGeofence(ctx context.Context, token string, projectID string, fence *lbs.CircleGeofence) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_circle_geofence for token %s and project %s and fanceName %s took %s to complete", token, projectID, fence.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateCircleGeofence(ctx, token, projectID, fence)
}

func (lm *loggingMiddleware) DeleteGeofence(ctx context.Context, token string, projectID string, fenceIDs []string, objects []string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method delete_geofence for token %s and project %s and fanceIds %s took %s to complete", token, projectID, fenceIDs, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.DeleteGeofence(ctx, token, projectID, fenceIDs, objects)
}

func (lm *loggingMiddleware) ListGeofences(ctx context.Context, token string, projectID string, fenceIDs []string, objects []string) (fenceList []*lbs.Geofence, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_geofences for token %s and project %s and fanceIds %s took %s to complete", token, projectID, fenceIDs, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListGeofences(ctx, token, projectID, fenceIDs, objects)
}

func (lm *loggingMiddleware) AddMonitoredObject(ctx context.Context, token string, projectID string, fenceID string, objects []string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_geofences for token %s and project %s and fanceId %s took %s to complete", token, projectID, fenceID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.AddMonitoredObject(ctx, token, projectID, fenceID, objects)
}

func (lm *loggingMiddleware) RemoveMonitoredObject(ctx context.Context, token string, projectID string, fenceID string, objects []string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_monitored_object for token %s and project %s and fanceId %s took %s to complete", token, projectID, fenceID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveMonitoredObject(ctx, token, projectID, fenceID, objects)
}

func (lm *loggingMiddleware) ListMonitoredObjects(ctx context.Context, token string, projectID string, fenceID string, pageIndex int, pageSize int) (total int, objects []string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_monitored_objects for token %s and project %s and fanceId %s and pageIndex %d and pageSize %d took %s to complete", token, projectID, fenceID, pageIndex, pageSize, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListMonitoredObjects(ctx, token, projectID, fenceID, pageIndex, pageSize)
}

func (lm *loggingMiddleware) CreatePolyGeofence(ctx context.Context, token string, projectID string, fence *lbs.PolyGeofence) (fenceID string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_poly_geofence for token %s and project %s and fanceName %s took %s to complete", token, projectID, fence.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreatePolyGeofence(ctx, token, projectID, fence)
}

func (lm *loggingMiddleware) UpdatePolyGeofence(ctx context.Context, token string, projectID string, fence *lbs.PolyGeofence) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_poly_geofence for token %s and project %s and fanceName %s took %s to complete", token, projectID, fence.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdatePolyGeofence(ctx, token, projectID, fence)
}

func (lm *loggingMiddleware) GetFenceIDs(ctx context.Context, token string, projectID string) (fenceIDs []string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method get_fenceids for token %s and project %s took %s to complete", token, projectID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.GetFenceIDs(ctx, token, projectID)
}

// Alarm
func (lm *loggingMiddleware) QueryStatus(ctx context.Context, token string, projectID string, monitoredPerson string, fenceIDs []string) (status *lbs.QueryStatus, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method query_status for token %s and project %s monitoredpersion %s fenceIDs %s took %s to complete", token, projectID, monitoredPerson, fenceIDs, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.QueryStatus(ctx, token, projectID, monitoredPerson, fenceIDs)
}

func (lm *loggingMiddleware) GetHistoryAlarms(ctx context.Context, token string, projectID string, monitoredPerson string, fenceIDs []string) (alarms *lbs.HistoryAlarms, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method get_history_alarms for token %s and project %s monitoredpersion %s fenceIDs %s took %s to complete", token, projectID, monitoredPerson, fenceIDs, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.GetHistoryAlarms(ctx, token, projectID, monitoredPerson, fenceIDs)
}

func (lm *loggingMiddleware) BatchGetHistoryAlarms(ctx context.Context, token string, projectID string, input *lbs.BatchGetHistoryAlarmsRequest) (alarms *lbs.BatchHistoryAlarmsResp, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method batchget_history_alarms for token %s and project %s took %s to complete", token, projectID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.BatchGetHistoryAlarms(ctx, token, projectID, input)
}

func (lm *loggingMiddleware) GetStayPoints(ctx context.Context, token string, projectID string, input *lbs.GetStayPointsRequest) (points *lbs.StayPoints, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method get_stay_points for token %s and project %s took %s to complete", token, projectID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.GetStayPoints(ctx, token, projectID, input)
}

// NotifyAlarms is used by apiserver to provide asynchrous notication
func (lm *loggingMiddleware) NotifyAlarms(ctx context.Context, token string, projectID string, content []byte) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method notify_alarms for token %s and project %s took %s to complete", token, projectID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.NotifyAlarms(ctx, token, projectID, content)
}

func (lm *loggingMiddleware) GetFenceUserID(ctx context.Context, token string, fenceID string) (userID string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method get_fence_userids for token %s and fenceID %s took %s to complete", token, fenceID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.GetFenceUserID(ctx, token, fenceID)
}

//Entity
func (lm *loggingMiddleware) AddEntity(ctx context.Context, token string, projectID string, entityName string, entityDesc string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_entiry for token %s and entityName %s took %s to complete", token, entityName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.AddEntity(ctx, token, projectID, entityName, entityDesc)
}

func (lm *loggingMiddleware) UpdateEntity(ctx context.Context, token string, projectID string, entityName string, entityDesc string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_entity for token %s and entityName %s took %s to complete", token, entityName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateEntity(ctx, token, projectID, entityName, entityDesc)
}

func (lm *loggingMiddleware) DeleteEntity(ctx context.Context, token string, projectID string, entityName string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method delete_entity for token %s and entityName %s took %s to complete", token, entityName, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.DeleteEntity(ctx, token, projectID, entityName)
}

func (lm *loggingMiddleware) ListEntity(ctx context.Context, token string, projectID string, coordTypeOutput string, pageIndex int, pageSize int) (total int, infos []*lbs.EntityInfo, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_entity for token %s and coordTypeOutput %s and pageIndex %d and pageSize %d took %s to complete", token, coordTypeOutput, pageIndex, pageSize, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListEntity(ctx, token, projectID, coordTypeOutput, pageIndex, pageSize)
}
