package lbs

import (
	"context"
	"errors"
	"fmt"

	nats "github.com/cloustone/pandas/lbs/nats/publisher"
	"github.com/cloustone/pandas/mainflux"
	logr "github.com/sirupsen/logrus"
)

var _ Service = (*lbsService)(nil)
var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrScanMetadata indicates problem with metadata in db
	ErrScanMetadata = errors.New("Failed to scan metadata")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	ListCollections(ctx context.Context, token string) ([]string, error)

	// Geofence
	CreateCircleGeofence(ctx context.Context, token string, projectId string, fence *CircleGeofence) (string, error)
	UpdateCircleGeofence(ctx context.Context, token string, projectId string, fence *CircleGeofence) error
	DeleteGeofence(ctx context.Context, token string, projectId string, fenceIDs []string, objects []string) error
	ListGeofences(ctx context.Context, token string, projectId string, fenceIDs []string, objects []string) ([]*Geofence, error)
	AddMonitoredObject(ctx context.Context, token string, projectId string, fenceID string, objects []string) error
	RemoveMonitoredObject(ctx context.Context, token string, projectId string, fenceID string, objects []string) error
	ListMonitoredObjects(ctx context.Context, token string, projectId string, fenceID string, pageIndex int, pageSize int) (int, []string, error)
	CreatePolyGeofence(ctx context.Context, token string, projectId string, fence *PolyGeofence) (string, error)
	UpdatePolyGeofence(ctx context.Context, token string, projectId string, fence *PolyGeofence) error
	GetFenceIDs(ctx context.Context, token string, projectId string) ([]string, error)

	// Alarm
	QueryStatus(ctx context.Context, token string, projectId string, monitoredPerson string, fenceIDs []string) (*QueryStatus, error)
	GetHistoryAlarms(ctx context.Context, token string, projectId string, monitoredPerson string, fenceIDs []string) (*HistoryAlarms, error)
	BatchGetHistoryAlarms(ct context.Context, token string, projectId string, input *BatchGetHistoryAlarmsRequest) (*BatchHistoryAlarmsResp, error)
	GetStayPoints(ctx context.Context, token string, projectId string, input *GetStayPointsRequest) (*StayPoints, error)

	// NotifyAlarms is used by apiserver to provide asynchrous notication
	NotifyAlarms(ctx context.Context, token string, projectId string, content []byte) error
	GetFenceUserID(ctx context.Context, token string, fenceID string) (string, error)

	//Entity
	AddEntity(ctx context.Context, token string, projectId string, entityName string, entityDesc string) error
	UpdateEntity(ctx context.Context, token string, projectId string, entityName string, entityDesc string) error
	DeleteEntity(ctx context.Context, token string, projectId string, entityName string) error

	ListEntity(ctx context.Context, token string, projectId string, coordTypeOutput string, pageIndex int, pageSize int) (int, []*EntityInfo, error)
}

// New instantiates the lbs service implementation.
func New(auth mainflux.AuthNServiceClient,
	provider LocationProvider, collections CollectionRepository, entities EntityRepository, geofences GeofenceRepository,
	n *nats.Publisher) Service {
	return &lbsService{
		auth:        auth,
		provider:    provider,
		collections: collections,
		entities:    entities,
		geofences:   geofences,
	}
}

type lbsService struct {
	auth        mainflux.AuthNServiceClient
	provider    LocationProvider
	collections CollectionRepository
	entities    EntityRepository
	geofences   GeofenceRepository
	nats        *nats.Publisher
}

// Geofence
func (l *lbsService) CreateCircleGeofence(ctx context.Context, token string, projectId string, fence *CircleGeofence) (string, error) {
	logr.Debugf("CreateCircleGeofence ()")

	res, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	userID := res.GetValue()

	name := fmt.Sprintf("%s-%s-%s", userID, projectId, fence.Name)
	fenceID, err := l.provider.CreateCircleGeofence(ctx,
		CircleGeofence{
			Name: name,
			//MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Longitude: fence.Longitude,
			Latitude:  fence.Latitude,
			Radius:    fence.Radius,
			Denoise:   int(fence.Denoise),
			CoordType: CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("create circle geofence failed")
		return "", errors.New("create circle geofence failed")
	}
	return fenceID, nil
}

func (l *lbsService) CreatePolyGeofence(ctx context.Context, token string, projectId string, fence *PolyGeofence) (string, error) {
	logr.Debugf("CreatePolyGeofence (%s)", fence.Name)

	res, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	userID := res.GetValue()

	name := fmt.Sprintf("%s-%s-%s", userID, projectId, fence.Name)
	fenceID, err := l.provider.CreatePolyGeofence(ctx,
		PolyGeofence{
			Name: name,
			//MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Vertexes:  fence.Vertexes,
			Denoise:   int(fence.Denoise),
			CoordType: CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("create poly geofence failed")
		return "", errors.New("create poly geofence failed")
	}
	return fenceID, nil
}

func (l *lbsService) UpdatePolyGeofence(ctx context.Context, token string, projectId string, fence *PolyGeofence) error {
	logr.Debugf("UpdatePolyGeofence (%s)", fence.Name)

	res, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	userID := res.GetValue()

	name := fmt.Sprintf("%s-%s-%s", userID, projectId, fence.Name)

	err = l.provider.UpdatePolyGeofence(ctx,
		PolyGeofence{
			Name: name,
			//			MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Vertexes:  fence.Vertexes,
			Denoise:   int(fence.Denoise),
			FenceID:   fence.FenceID,
			CoordType: CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("update poly geofence failed")
		return errors.New("update poly geofence failed")
	}
	return nil
}

func (l *lbsService) UpdateCircleGeofence(ctx context.Context, token string, projectId string, fence *CircleGeofence) error {
	logr.Debugf("UpdateCircleGeofence ()")

	res, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	userID := res.GetValue()

	name := fmt.Sprintf("%s-%s-%s", userID, projectId, fence.Name)

	err = l.provider.UpdateCircleGeofence(ctx,
		CircleGeofence{
			Name: name,
			//		MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Longitude: fence.Longitude,
			Latitude:  fence.Latitude,
			Radius:    fence.Radius,
			Denoise:   int(fence.Denoise),
			FenceID:   fence.FenceID,
			CoordType: CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("update circle geofence failed")
		return errors.New("update circle geofence failed")
	}
	return nil
}

func (l *lbsService) DeleteGeofence(ctx context.Context, token string, projectId string, fenceIDs []string, objects []string) error {
	logr.Debugf("DeleteGeofence (%s)", fenceIDs)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	_, err = l.provider.DeleteGeofence(ctx,
		fenceIDs, objects)
	if err != nil {
		logr.WithError(err).Errorf("delete geofence failed")
		return errors.New("delete geofence failed")
	}
	return nil
}

func (l *lbsService) ListGeofences(ctx context.Context, token string, projectId string, fenceIDs []string, objects []string) ([]*Geofence, error) {
	logr.Debugf("ListGeofence (%s)", fenceIDs)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	fences, err := l.provider.ListGeofence(ctx,
		fenceIDs, objects)
	if err != nil {
		logr.WithError(err).Errorf("list geofence failed")
		return nil, errors.New("list geofence failed")
	}
	fenceList := []*Geofence{}
	for _, f := range fences {
		fence := &Geofence{
			//		FenceID:         fmt.Sprint(f.FenceID),
			FenceName: f.FenceName,
			//MonitoredObject: strings.Split(f.MonitoredObject, ","),
			Shape:     f.Shape,
			Longitude: f.Longitude,
			Latitude:  f.Latitude,
			Radius:    f.Radius,
			CoordType: f.CoordType,
			//Denoise:    int(f.Denoise),
			//CreateTime: f.CreateTime,
			UpdateTime: f.UpdateTime,
		}
		for _, vtx := range f.Vertexes {
			vertexe := Vertexe{
				Latitude:  vtx.Latitude,
				Longitude: vtx.Longitude,
			}
			fence.Vertexes = append(fence.Vertexes, vertexe)
		}
		fenceList = append(fenceList, fence)
	}
	return fenceList, nil
}

func (l *lbsService) AddMonitoredObject(ctx context.Context, token string, projectId string, fenceID string, objects []string) error {
	logr.Debugf("AddMonitoredObject (%s)", fenceID)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.provider.AddMonitoredObject(ctx,
		fenceID, objects)
	if err != nil {
		logr.WithError(err).Errorf("add monitored object failed")
		return errors.New("ladd monitored object failed")
	}
	return nil
}

func (l *lbsService) RemoveMonitoredObject(ctx context.Context, token string, projectId string, fenceID string, objects []string) error {
	logr.Debugf("RemoveMonitoredObject(%s)", fenceID)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.provider.RemoveMonitoredObject(ctx,
		fenceID, objects)
	if err != nil {
		logr.WithError(err).Errorf("remove monitored object failed")
		return errors.New("remove monitored object failed")
	}
	return nil
}

func (l *lbsService) ListMonitoredObjects(ctx context.Context, token string, projectId string, fenceID string, pageIndex int, pageSize int) (int, []string, error) {
	logr.Debugf("ListMonitoredObjects(%s)", fenceID)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return 0, nil, ErrUnauthorizedAccess
	}

	total, objects := l.provider.ListMonitoredObjects(ctx,
		fenceID, int(pageIndex), int(pageSize))

	return int(total), objects, nil
}

func (l *lbsService) ListCollections(ctx context.Context, token string) ([]string, error) {
	logr.Debugln("ListCollections")

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	//	productList := getProductList(products)
	//	logr.Debugln("productList:", productList)

	//	return productList, nil
	return nil, nil
}

func (l *lbsService) GetFenceIDs(ctx context.Context, token string, projectId string) ([]string, error) {
	logr.Debugf("GetFenceIDs (%s)", projectId)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}
	/*

		fences, err := l.provider.GetFenceIDs(ctx,
			auth.NewPrincipal(in.UserID, in.ProjectID))
		if err != nil {
			logr.WithError(err).Errorln("ListCollections err:", err)
			return nil, err
		}

		fenceIDs := getFenceIDs(fences)

		return fenceIDs, nil
	*/
	return nil, nil
}

func (l *lbsService) QueryStatus(ctx context.Context, token string, projectId string, monitoredPerson string, fenceIDs []string) (*QueryStatus, error) {
	logr.Debugf("QueryStatus(%s)", fenceIDs)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	fenceStatus, err := l.provider.QueryStatus(ctx,
		monitoredPerson, fenceIDs)
	if err != nil {
		logr.Errorln("QueryStatus failed:", err)
		return nil, errors.New("not found")
	}

	return &fenceStatus, nil
}

func (l *lbsService) GetHistoryAlarms(ctx context.Context, token string, projectId string, monitoredPerson string, fenceIDs []string) (*HistoryAlarms, error) {
	logr.Debugf("GetHistoryAlarms (%s)", fenceIDs)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	alarmPoint, err := l.provider.GetHistoryAlarms(ctx,
		monitoredPerson, fenceIDs)
	if err != nil {
		logr.Errorln("GetHistoryAlarms failed:", err)
		return nil, errors.New("not found")
	}

	return &alarmPoint, nil
}

func (l *lbsService) BatchGetHistoryAlarms(ctx context.Context, token string, projectId string, input *BatchGetHistoryAlarmsRequest) (*BatchHistoryAlarmsResp, error) {
	logr.Debugf("RemoveCollection ")

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	historyAlarms, err := l.provider.BatchGetHistoryAlarms(ctx,
		input)
	if err != nil {
		logr.Errorln("BatchGetHistoryAlarms failed:", err)
		return nil, errors.New("not found")
	}

	return &historyAlarms, nil

}

func (l *lbsService) GetStayPoints(ctx context.Context, token string, projectId string, input *GetStayPointsRequest) (*StayPoints, error) {
	logr.Debugf("GetStayPoints ()")

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	stayPoints, err := l.provider.GetStayPoints(ctx,
		input)
	if err != nil {
		logr.Errorln("BatchGetHistoryAlarms failed:", err)
		return nil, errors.New("not found")
	}

	return &stayPoints, nil
}

func (l *lbsService) NotifyAlarms(ctx context.Context, token string, projectId string, content []byte) error {
	logr.Debugf("NotifyAlarms(%s)", content)

	/*
		_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
		if err != nil {
			return ErrUnauthorizedAccess
		}

		alarm, err := l.provider.HandleAlarmNotification(ctx,
			content)
		if err != nil {
			logr.WithError(err).Errorf("unmarshal alarm failed")
			return errors.New("unmarshal alarm failed")
		}
		config := message.NewConfigWithViper()
		producer, err := message.NewProducer(config, false)
		if err != nil {
			logr.WithError(err).Errorf("create message producer failed")
			return errors.New("create message producer failed")
		}
			if err := producer.SendMessage(&AlarmTopic{Alarm: alarm}); err != nil {
				logr.WithError(err).Errorf("send alarm failed")
				return errors.New("send alarm failed")
			}
	*/
	return nil
}

func (l *lbsService) GetFenceUserID(ctx context.Context, token string, fenceID string) (string, error) {
	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	/*
		userID, err := l.provider.GetFenceUserID(ctx, nil, fenceID)
		if err != nil {
			logr.WithError(err).Errorf("add collection '%s' failed", fenceID)
			return nil, err
		}
		return userID, nil
	*/
	return "", nil
}

/*
func getProductList(products []*Collection) *pb.ListCollectionsResponse {
	productList := &pb.ListCollectionsResponse{}
	for _, val := range products {
		productList.ProjectIDs = append(productList.ProjectIDs, val.ProjectID)
	}
	return productList
}
*/

func (l *lbsService) AddEntity(ctx context.Context, token string, projectId string, entityName string, entityDesc string) error {
	logr.Debugf("AddEntity (%s)", entityName)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.provider.AddEntity(ctx,
		entityName, entityDesc)
	if err != nil {
		logr.WithError(err).Errorf("AddEntity failed")
		return errors.New("AddEntity failed")
	}
	return nil
}

func (l *lbsService) DeleteEntity(ctx context.Context, token string, projectId string, entityName string) error {
	logr.Debugf("DeleteEntity (%s)", entityName)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.provider.DeleteEntity(ctx,
		entityName)
	if err != nil {
		logr.WithError(err).Errorf("DeleteEntity failed")
		return errors.New("DeleteEntity failed")
	}
	return nil
}

func (l *lbsService) UpdateEntity(ctx context.Context, token string, projectId string, entityName string, entityDesc string) error {
	logr.Debugf("UpdateEntity (%s)", entityName)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.provider.UpdateEntity(ctx,
		entityName, entityDesc)
	if err != nil {
		logr.WithError(err).Errorf("UpdateEntity failed")
		return errors.New("UpdateEntity failed")
	}
	/*
		entity := EntityRecord{
			//	UserID:        in.UserID,
			//	ProjectID:     in.ProjectID,
			EntityName:    in.EntityName,
			LastUpdatedAt: time.Now(),
		}
	*/

	return nil
}

func (l *lbsService) ListEntity(ctx context.Context, token string, projectId string, coordTypeOutput string, pageIndex int, pageSize int) (int, []*EntityInfo, error) {
	logr.Debugf("ListEntity (%s)", coordTypeOutput)
	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return 0, nil, ErrUnauthorizedAccess
	}
	/*
		entitiesInfo, err := l.provider.GetEntity(ctx,
			auth.NewPrincipal(userID, projectId))
		if err != nil {
			logr.WithError(err).Errorf("mongo get entities info failed")
			return nil, nil, errors.Wrap(errors.New("mongo get entities info faile"), err)
		}
		entitiesName := getEntitiesName(entitiesInfo)

		total, entityInfo := l.provider.ListEntity(ctx, userID, projectId, coordTypeOutput, pageIndex, pageSize)
		if total == -1 {
			logr.Errorf("ListEntity failed")
			return nil, nil, errors.Wrap(errors.New("ListEntity failed"), err)
		}

		entitys := &EntityInfo
		for _, val := range entityInfo {
			entityInfo := &EntityInfo{
				EntityName: val.EntityName,
				Longitude:  val.LastLocation.Longitude,
				Latitude:   val.LastLocation.Latitude,
			}
			if isEntityInCollection(val.EntityName, entitiesName) == false {
				continue
			}
			entitys = append(entitys, entityInfo)
		}
		Total = int(len(entitys))

		logr.Debugln("EntityInfo:", entityInfo, "total:", total)

		return Total,entitys, nil
	*/
	return 0, nil, nil
}

func getEntitiesName(entitiesInfo []*EntityRecord) []string {
	entitiesName := make([]string, 0)
	for _, val := range entitiesInfo {
		entitiesName = append(entitiesName, val.EntityName)
	}
	return entitiesName
}
func isEntityInCollection(entityName string, entitiesName []string) bool {
	var isEntityExist bool = false
	for _, val := range entitiesName {
		if val == entityName {
			isEntityExist = true
		}
	}
	return isEntityExist
}

/*
func getFenceIDs(fences []*GeofenceRecord) *pb.GetFenceIDsResponse {
	fenceIDsResp := &pb.GetFenceIDsResponse{}
	for _, val := range fences {
		fenceIDsResp.FenceIDs = append(fenceIDsResp.FenceIDs, val.FenceID)
	}
	logr.Debugln("fenceIDs:", fenceIDsResp.FenceIDs)
	return fenceIDsResp
}
*/
