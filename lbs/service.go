package lbs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	lbp "github.com/cloustone/pandas/lbs/proxy"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/message"
	logr "github.com/sirupsen/logrus"
)

var _ Service = (*LbsService)(nil)
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
	DeleteGeofence(ctx context.Context, token string, projectId string, fenceIds []string, objects []string) error
	ListGeofences(ctx context.Context, token string, projectId string, fenceIds []string, objects []string) ([]*Geofence, error)
	AddMonitoredObject(ctx context.Context, token string, projectId string, fenceId string, objects []string) error
	RemoveMonitoredObject(ctx context.Context, token string, projectId string, fenceId string, objects []string) error
	ListMonitoredObjects(ctx context.Context, token string, projectId string, fenceId string, pageIndex int32, pageSize int32) (int32, []string, error)
	CreatePolyGeofence(ctx context.Context, token string, projectId string, fence *PolyGeofence) (string, error)
	UpdatePolyGeofence(ctx context.Context, token string, projectId string, fence *PolyGeofence) error
	GetFenceIds(ctx context.Context, token string, projectId string) ([]string, error)

	// Alarm
	QueryStatus(ctx context.Context, token string, projectId string, monitoredPerson string, fenceIds []string) (*lbp.QueryStatus, error)
	GetHistoryAlarms(ctx context.Context, token string, projectId string, monitoredPerson string, fenceIds []string) (*lbp.HistoryAlarms, error)
	BatchGetHistoryAlarms(ct context.Context, token string, projectId string, input *lbp.BatchGetHistoryAlarmsRequest) (*lbp.HistoryAlarms, error)
	GetStayPoints(ctx context.Context, token string, projectId string, input *lbp.GetStayPointsRequest) (*lbp.StayPoints, error)

	// NotifyAlarms is used by apiserver to provide asynchrous notication
	NotifyAlarms(ctx context.Context, token string, projectId string, content []byte) error
	GetFenceUserId(ctx context.Context, token string, fenceId string) (string, error)

	//Entity
	AddEntity(ctx context.Context, token string, projectId string, entityName string, entityDesc string) error
	UpdateEntity(ctx context.Context, token string, projectId string, entityName string, entityDesc string) error
	DeleteEntity(ctx context.Context, token string, projectId string, entityName string) error

	ListEntity(ctx context.Context, token string, projectId string, coordTypeOutput string, pageIndex int32, pageSize int32) (int32, []*EntityInfo, error)
}

// Geofence
type CircleGeofence struct {
	Name             string
	MonitoredObjects []string
	Longitude        float64
	Latitude         float64
	Radius           float64
	CoordType        string
	Denoise          int32
	FenceId          string
}

type Geofence struct {
	FenceId         string
	FenceName       string
	MonitoredObject []string
	Shape           string
	Longitude       float64
	Latitude        float64
	Radius          float64
	CoordType       string
	Denoise         int32
	CreateTime      string
	UpdateTime      string
	Vertexes        []*Vertexe
}
type Vertexe struct {
	Longitude float64
	Latitude  float64
}

type PolyGeofence struct {
	Name             string
	MonitoredObjects []string
	Vertexes         string
	CoordType        string
	Denoise          int32
	FenceId          string
}

type LbsService struct {
	auth  mainflux.AuthNServiceClient
	Proxy lbp.Proxy
}

// Alarm
type MonitoredStatus struct {
	FenceId         int32
	MonitoredStatus string
}

type AlarmPoint struct {
	Longitude  float64
	Latitude   float64
	Radius     int32
	CoordType  string
	LocTime    string
	CreateTime string
}

type PrePoint struct {
	Longitude  float64
	Latitude   float64
	Radius     int32
	CoordType  string
	LocTime    string
	CreateTime string
}

type Alarm struct {
	FenceId         int32
	FenceName       string
	MonitoredPerson string
	Action          string
	AlarmPoint      *AlarmPoint
	PrePoint        *PrePoint
}

type EntityInfo struct {
	EntityName string
	Latitude   float64
	Longitude  float64
}

// New instantiates the lbs service implementation.
func New(proxy lbp.Proxy) Service {
	return &LbsService{
		Proxy: proxy,
	}
}

// Geofence
func (l *LbsService) CreateCircleGeofence(ctx context.Context, token string, projectId string, fence *CircleGeofence) (string, error) {
	logr.Debugf("CreateCircleGeofence ()")

	res, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	userId := res.GetValue()

	name := fmt.Sprintf("%s-%s-%s", userId, projectId, fence.Name)
	fenceId, err := l.Proxy.CreateCircleGeofence(
		lbp.CircleGeofence{
			Name:             name,
			MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Longitude:        fence.Longitude,
			Latitude:         fence.Latitude,
			Radius:           fence.Radius,
			Denoise:          int(fence.Denoise),
			CoordType:        lbp.CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("create circle geofence failed")
		return "", errors.New("create circle geofence failed")
	}
	return fenceId, nil
}

func (l *LbsService) CreatePolyGeofence(ctx context.Context, token string, projectId string, fence *PolyGeofence) (string, error) {
	logr.Debugf("CreatePolyGeofence (%s)", fence.Name)

	res, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	userId := res.GetValue()

	name := fmt.Sprintf("%s-%s-%s", userId, projectId, fence.Name)
	fenceId, err := l.Proxy.CreatePolyGeofence(
		lbp.PolyGeofence{
			Name:             name,
			MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Vertexes:         fence.Vertexes,
			Denoise:          int(fence.Denoise),
			CoordType:        lbp.CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("create poly geofence failed")
		return "", errors.New("create poly geofence failed")
	}
	return fenceId, nil
}

func (l *LbsService) UpdatePolyGeofence(ctx context.Context, token string, projectId string, fence *PolyGeofence) error {
	logr.Debugf("UpdatePolyGeofence (%s)", fence.Name)

	res, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	userId := res.GetValue()

	name := fmt.Sprintf("%s-%s-%s", userId, projectId, fence.Name)

	err = l.Proxy.UpdatePolyGeofence(
		lbp.PolyGeofence{
			Name:             name,
			MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Vertexes:         fence.Vertexes,
			Denoise:          int(fence.Denoise),
			FenceId:          fence.FenceId,
			CoordType:        lbp.CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("update poly geofence failed")
		return errors.New("update poly geofence failed")
	}
	return nil
}

func (l *LbsService) UpdateCircleGeofence(ctx context.Context, token string, projectId string, fence *CircleGeofence) error {
	logr.Debugf("UpdateCircleGeofence ()")

	res, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}
	userId := res.GetValue()

	name := fmt.Sprintf("%s-%s-%s", userId, projectId, fence.Name)

	err = l.Proxy.UpdateCircleGeofence(
		lbp.CircleGeofence{
			Name:             name,
			MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Longitude:        fence.Longitude,
			Latitude:         fence.Latitude,
			Radius:           fence.Radius,
			Denoise:          int(fence.Denoise),
			FenceId:          fence.FenceId,
			CoordType:        lbp.CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("update circle geofence failed")
		return errors.New("update circle geofence failed")
	}
	return nil
}

func (l *LbsService) DeleteGeofence(ctx context.Context, token string, projectId string, fenceIds []string, objects []string) error {
	logr.Debugf("DeleteGeofence (%s)", fenceIds)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	_, err = l.Proxy.DeleteGeofence(
		fenceIds, objects)
	if err != nil {
		logr.WithError(err).Errorf("delete geofence failed")
		return errors.New("delete geofence failed")
	}
	return nil
}

func (l *LbsService) ListGeofences(ctx context.Context, token string, projectId string, fenceIds []string, objects []string) ([]*Geofence, error) {
	logr.Debugf("ListGeofence (%s)", fenceIds)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	fences, err := l.Proxy.ListGeofence(
		fenceIds, objects)
	if err != nil {
		logr.WithError(err).Errorf("list geofence failed")
		return nil, errors.New("list geofence failed")
	}
	fenceList := []*Geofence{}
	for _, f := range fences {
		fence := &Geofence{
			FenceId:         fmt.Sprint(f.FenceId),
			FenceName:       f.FenceName,
			MonitoredObject: strings.Split(f.MonitoredObject, ","),
			Shape:           f.Shape,
			Longitude:       f.Longitude,
			Latitude:        f.Latitude,
			Radius:          f.Radius,
			CoordType:       string(f.CoordType),
			Denoise:         int32(f.Denoise),
			CreateTime:      f.CreateTime,
			UpdateTime:      f.UpdateTime,
		}
		for _, vtx := range f.Vertexes {
			vertexe := &Vertexe{
				Latitude:  vtx.Latitude,
				Longitude: vtx.Longitude,
			}
			fence.Vertexes = append(fence.Vertexes, vertexe)
		}
		fenceList = append(fenceList, fence)
	}
	return fenceList, nil
}

func (l *LbsService) AddMonitoredObject(ctx context.Context, token string, projectId string, fenceId string, objects []string) error {
	logr.Debugf("AddMonitoredObject (%s)", fenceId)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.Proxy.AddMonitoredObject(
		fenceId, objects)
	if err != nil {
		logr.WithError(err).Errorf("add monitored object failed")
		return errors.New("ladd monitored object failed")
	}
	return nil
}

func (l *LbsService) RemoveMonitoredObject(ctx context.Context, token string, projectId string, fenceId string, objects []string) error {
	logr.Debugf("RemoveMonitoredObject(%s)", fenceId)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.Proxy.RemoveMonitoredObject(
		fenceId, objects)
	if err != nil {
		logr.WithError(err).Errorf("remove monitored object failed")
		return errors.New("remove monitored object failed")
	}
	return nil
}

func (l *LbsService) ListMonitoredObjects(ctx context.Context, token string, projectId string, fenceId string, pageIndex int32, pageSize int32) (int32, []string, error) {
	logr.Debugf("ListMonitoredObjects(%s)", fenceId)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return 0, nil, ErrUnauthorizedAccess
	}

	total, objects := l.Proxy.ListMonitoredObjects(
		fenceId, int(pageIndex), int(pageSize))

	return int32(total), objects, nil
}

func (l *LbsService) ListCollections(ctx context.Context, token string) ([]string, error) {
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

func (l *LbsService) GetFenceIds(ctx context.Context, token string, projectId string) ([]string, error) {
	logr.Debugf("GetFenceIds (%s)", projectId)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}
	/*

		fences, err := l.Proxy.GetFenceIds(
			auth.NewPrincipal(in.UserId, in.ProjectId))
		if err != nil {
			logr.WithError(err).Errorln("ListCollections err:", err)
			return nil, err
		}

		fenceIds := getFenceIds(fences)

		return fenceIds, nil
	*/
	return nil, nil
}

func (l *LbsService) QueryStatus(ctx context.Context, token string, projectId string, monitoredPerson string, fenceIds []string) (*lbp.QueryStatus, error) {
	logr.Debugf("QueryStatus(%s)", fenceIds)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	fenceStatus, err := l.Proxy.QueryStatus(
		monitoredPerson, fenceIds)
	if err != nil {
		logr.Errorln("QueryStatus failed:", err)
		return nil, errors.New("not found")
	}

	return fenceStatus, nil
}

func (l *LbsService) GetHistoryAlarms(ctx context.Context, token string, projectId string, monitoredPerson string, fenceIds []string) (*lbp.HistoryAlarms, error) {
	logr.Debugf("GetHistoryAlarms (%s)", fenceIds)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	alarmPoint, err := l.Proxy.GetHistoryAlarms(
		monitoredPerson, fenceIds)
	if err != nil {
		logr.Errorln("GetHistoryAlarms failed:", err)
		return nil, errors.New("not found")
	}

	return alarmPoint, nil
}

func (l *LbsService) BatchGetHistoryAlarms(ctx context.Context, token string, projectId string, input *lbp.BatchGetHistoryAlarmsRequest) (*lbp.HistoryAlarms, error) {
	logr.Debugf("RemoveCollection ")

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	historyAlarms, err := l.Proxy.BatchGetHistoryAlarms(
		input)
	if err != nil {
		logr.Errorln("BatchGetHistoryAlarms failed:", err)
		return nil, errors.New("not found")
	}

	return historyAlarms, nil

}

func (l *LbsService) GetStayPoints(ctx context.Context, token string, projectId string, input *lbp.GetStayPointsRequest) (*lbp.StayPoints, error) {
	logr.Debugf("GetStayPoints ()")

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	stayPoints, err := l.Proxy.GetStayPoints(
		input)
	if err != nil {
		logr.Errorln("BatchGetHistoryAlarms failed:", err)
		return nil, errors.New("not found")
	}

	return stayPoints, nil
}

func (l *LbsService) NotifyAlarms(ctx context.Context, token string, projectId string, content []byte) error {
	logr.Debugf("NotifyAlarms(%s)", content)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	alarm, err := l.Proxy.UnmarshalAlarmNotification(
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

	if err := producer.SendMessage(&lbp.AlarmTopic{Alarm: alarm}); err != nil {
		logr.WithError(err).Errorf("send alarm failed")
		return errors.New("send alarm failed")
	}
	return nil
}

func (l *LbsService) GetFenceUserId(ctx context.Context, token string, fenceId string) (string, error) {
	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	/*
		userId, err := l.Proxy.GetFenceUserId(nil, fenceId)
		if err != nil {
			logr.WithError(err).Errorf("add collection '%s' failed", fenceId)
			return nil, err
		}
		return userId, nil
	*/
	return "", nil
}

/*
func getProductList(products []*Collection) *pb.ListCollectionsResponse {
	productList := &pb.ListCollectionsResponse{}
	for _, val := range products {
		productList.ProjectIds = append(productList.ProjectIds, val.ProjectId)
	}
	return productList
}
*/

func (l *LbsService) AddEntity(ctx context.Context, token string, projectId string, entityName string, entityDesc string) error {
	logr.Debugf("AddEntity (%s)", entityName)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.Proxy.AddEntity(
		entityName, entityDesc)
	if err != nil {
		logr.WithError(err).Errorf("AddEntity failed")
		return errors.New("AddEntity failed")
	}
	return nil
}

func (l *LbsService) DeleteEntity(ctx context.Context, token string, projectId string, entityName string) error {
	logr.Debugf("DeleteEntity (%s)", entityName)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.Proxy.DeleteEntity(
		entityName)
	if err != nil {
		logr.WithError(err).Errorf("DeleteEntity failed")
		return errors.New("DeleteEntity failed")
	}
	return nil
}

func (l *LbsService) UpdateEntity(ctx context.Context, token string, projectId string, entityName string, entityDesc string) error {
	logr.Debugf("UpdateEntity (%s)", entityName)

	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	err = l.Proxy.UpdateEntity(
		entityName, entityDesc)
	if err != nil {
		logr.WithError(err).Errorf("UpdateEntity failed")
		return errors.New("UpdateEntity failed")
	}
	/*
		entity := EntityRecord{
			//	UserId:        in.UserId,
			//	ProjectId:     in.ProjectId,
			EntityName:    in.EntityName,
			LastUpdatedAt: time.Now(),
		}
	*/

	return nil
}

func (l *LbsService) ListEntity(ctx context.Context, token string, projectId string, coordTypeOutput string, pageIndex int32, pageSize int32) (int32, []*EntityInfo, error) {
	logr.Debugf("ListEntity (%s)", coordTypeOutput)
	_, err := l.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return 0, nil, ErrUnauthorizedAccess
	}
	/*
		entitiesInfo, err := l.Proxy.GetEntity(
			auth.NewPrincipal(userId, projectId))
		if err != nil {
			logr.WithError(err).Errorf("mongo get entities info failed")
			return nil, nil, errors.Wrap(errors.New("mongo get entities info faile"), err)
		}
		entitiesName := getEntitiesName(entitiesInfo)

		total, entityInfo := l.Proxy.ListEntity(userId, projectId, coordTypeOutput, pageIndex, pageSize)
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
		Total = int32(len(entitys))

		logr.Debugln("EntityInfo:", entityInfo, "total:", total)

		return Total,entitys, nil
	*/
	return 0, nil, nil
}

func getEntitiesName(entitiesInfo []*lbp.EntityRecord) []string {
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
func getFenceIds(fences []*GeofenceRecord) *pb.GetFenceIdsResponse {
	fenceIdsResp := &pb.GetFenceIdsResponse{}
	for _, val := range fences {
		fenceIdsResp.FenceIds = append(fenceIdsResp.FenceIds, val.FenceId)
	}
	logr.Debugln("fenceIds:", fenceIdsResp.FenceIds)
	return fenceIdsResp
}
*/
