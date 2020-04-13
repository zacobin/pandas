package lbs

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/cloustone/pandas/lbs/grpc_lbs_v1"
	lbp "github.com/cloustone/pandas/lbs/proxy"
	"github.com/cloustone/pandas/pkg/auth"
	"github.com/cloustone/pandas/pkg/message"
	logr "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/cloustone/pandas/pkg/errors"
)

var gerrf = status.Errorf

var _ Service = (*LbsService)(nil)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	ListCollections(ctx context.Context, userId string) ([]string, error)

	// Geofence
	CreateCircleGeofence(ctx context.Context, userId string, projectId string, fence *CircleGeofence) (string, error)
	UpdateCircleGeofence(ctx context.Context, userId string, projectId string, fence *CircleGeofence) error
	DeleteGeofence(ctx context.Context, userId string, projectId string, fenceIds []string], objects []stirng) ([]string, error)
	ListGeofences(ctx context.Context, userId string, projectId string, fenceIds []string], objects []stirng) ([]*Geofence, error)
	AddMonitoredObject(ctx context.Context, userId string, projectId string, fenceId string, objects []stirng) error
	RemoveMonitoredObject(ctx context.Context, userId string, projectId string, fenceId string, objects []stirng) error
	ListMonitoredObjects(ctx context.Context, userId string, projectId string, fenceId string, pageIndex int32, pageSize int32) (int32, []string, error)
	CreatePolyGeofence(ctx context.Context, userId string, projectId string, fence *PolyGeofence) (string, error)
	UpdatePolyGeofence(ctx context.Context, userId string, projectId string, fence *PolyGeofence) error
	GetFenceIds(ctx context.Context, userId string, projectId string, fenceIds []string) ([]string], error)

	// Alarm
	QueryStatus(ctx context.Context, userId string, projectId string, monitoredPerson string, fendeIds []string) (*QueryStatus, error)
	GetHistoryAlarms(ctx context.Context, userId string, projectId string, monitoredPerson string, fendeIds []string) (*HistoryAlarms, error)
	BatchGetHistoryAlarms(ct context.Context, userId string, projectId string, coordTypeOutput string, endTime string, startTime string, pageIndex int32, pageSize int32) (*HistoryAlarms, error)
	GetStayPoints(ctx context.Context, userId string, projectId string, coordTypeOutput string, endTime string, startTime string, pageIndex int32, pageSize int32, fenceIds []string, entityName string) (*StayPoints, error) 

	// NotifyAlarms is used by apiserver to provide asynchrous notication
	NotifyAlarms(ctx context.Context, userId string, projectId string, content string) error
	GetFenceUserId(ctx context.Context, fenceId string) (string, error)

	//Entity
	AddEntity(ctx context.Context, userId string, projectId string, entityName string, entityDesc string) error
	UpdateEntity(ctx context.Context, userId string, projectId string, entityName string, entityDesc string) error
	DeleteEntity(ctx context.Context, userId string, projectId string, entityName string) error
   
   
    ListEntity(ListEntityRequest) returns (ListEntityResponse) {}

}

// Geofence
type CircleGeofence struct {
	Name                 string
	MonitoredObjects     []string
	Longitude            float64
	Latitude             float64
	Radius               float64
	CoordType            string
	Denoise              int32
	FenceId              string
}

type Geofence struct {
	FenceId              string     
	FenceName            string     
	MonitoredObject      []string   
	Shape                string     
	Longitude            float64    
	Latitude             float64    
	Radius               float64    
	CoordType            string     
	Denoise              int32      
	CreateTime           string     
	UpdateTime           string     
	Vertexes             []*Vertexe 
}
type Vertexe struct {
	Longitude            float64   
	Latitude             float64  
}

type PolyGeofence struct {
	Name                 string   
	MonitoredObjects     []string 
	Vertexes             string  
	CoordType            string   
	Denoise              int32    
	FenceId              string   
}

type LbsService struct {
	Proxy *lbp.Proxy
}

// Alarm
type MonitoredStatus struct {
	FenceId              int32    
	MonitoredStatus      string   
}

type QueryStatus struct {
	Status               int32              
	Message              string             
	Size                 int32              
	MonitoredStatuses    []*MonitoredStatus 
}

type AlarmPoint struct {
	Longitude            float64  
	Latitude             float64  
	Radius               int32    
	CoordType            string   
	LocTime              string   
	CreateTime           string   
}

type PrePoint struct {
	Longitude            float64  
	Latitude             float64  
	Radius               int32    
	CoordType            string   
	LocTime              string   
	CreateTime           string   
}
type Point struct {
	Latitude             float64  
	Longitude            float64  
	CoordType            string   
	LocTime              string   
}

type Alarm struct {
	FenceId              int32       
	FenceName            string      
	MonitoredPerson      string      
	Action               string      
	AlarmPoint           *AlarmPoint 
	PrePoint             *PrePoint   
}

type HistoryAlarms struct {
	Status               int32    
	Message              string   
	Total                int32
	Size                 int32    
	Alarms               []*Alarm 
}
type StayPoints struct {
	Status               int32    
	Message              string   
	Total                int32    
	Size                 int32    
	Distance             int32    
	EndPoint             *Point   
	StartPoint           *Point   
	Points               []*Point 
}

type EntityInfo struct {
	EntityName           string   
	Latitude             float64  
	Longitude            float64  
}

// New instantiates the lbs service implementation.
func New(proxy *lbp.Proxy) Service {
	return &service{
		Proxy: proxy,
	}
}

// Geofence
func (l *LbsService) CreateCircleGeofence(ctx context.Context, userId string, projectId string, fence *CircleGeofence) (string, error) {
	logr.Debugf("CreateCircleGeofence (%s)", in.String())

	name := fmt.Sprintf("%s-%s-%s", userId, projectId, fence.Name)
	fenceId, err := l.Proxy.CreateCircleGeofence(
		auth.NewPrincipal(userId, projectId),
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
		return nil, errors.Wrap("create circle geofence failed", err)
	}
	return fenceId, nil
}

type CreatePolyGeofenceRequest struct {
	UserId               string        `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty", bson:"user_id,omitempty"`
	ProjectId            string        `protobuf:"bytes,2,opt,name=project_id,json=projectId" json:"project_id,omitempty", bson:"project_id,omitempty"`
	Fence                *PolyGeofence `protobuf:"bytes,3,opt,name=fence" json:"fence,omitempty", bson:"fence,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}
type CreatePolyGeofenceResponse struct {
	FenceId              string   `protobuf:"bytes,1,opt,name=fence_id,json=fenceId" json:"fence_id,omitempty", bson:"fence_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (l *LbsService) CreatePolyGeofence(ctx context.Context, userId string, projectId string, fence *PolyGeofence) (string, error) {
	logr.Debugf("CreatePolyGeofence (%s)", fence.Name)

	name := fmt.Sprintf("%s-%s-%s", userId, projectId, fence.Name)
	fenceId, err := l.Proxy.CreatePolyGeofence(
		auth.NewPrincipal(userId, projectId),
		lbp.PolyGeofence{
			Name:             name,
			MonitoredObjects: strings.Join(fence.MonitoredObjects, ","),
			Vertexes:         fence.Vertexes,
			Denoise:          int(fence.Denoise),
			CoordType:        lbp.CoordType(fence.CoordType),
		})
	if err != nil {
		logr.WithError(err).Errorf("create poly geofence failed")
		return nil, errors.Wrap("create poly geofence failed", err)
	}
	return fenceId, nil
}

func (l *LbsService) UpdatePolyGeofence(ctx context.Context, userId string, projectId string, fence *PolyGeofence) error {
	logr.Debugf("UpdatePolyGeofence (%s)", fence.Name)

	name := fmt.Sprintf("%s-%s-%s", userId, projectId, fence.Name)

	err := l.Proxy.UpdatePolyGeofence(
		auth.NewPrincipal(userId, projectId),
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
		return errors.Wrap("update poly geofence failed", err)
	}
	return nil
}

func (l *LbsService) UpdateCircleGeofence(ctx context.Context, userId string, projectId string, fence *CircleGeofence) error {
	logr.Debugf("UpdateCircleGeofence (%s)", in.String())

	name := fmt.Sprintf("%s-%s-%s", userId, projectId, fence.Name)

	err := l.Proxy.UpdateCircleGeofence(
		auth.NewPrincipal(userId, projectId),
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
		return errors.Wrap("update circle geofence failed", err)
	}
	return nil
}

func (l *LbsService) DeleteGeofence(ctx context.Context, userId string, projectId string, fenceIds []string], objects []stirng) ([]string, error) {
	logr.Debugf("DeleteGeofence (%s)", fenceIds)

	fIds, err := l.Proxy.DeleteGeofence(
		auth.NewPrincipal(userId, projectId),
		fenceIds, objects)
	if err != nil {
		logr.WithError(err).Errorf("delete geofence failed")
		return nil, errors.Wrap("delete geofence failed", err)
	}
	return fIds, nil
}

func (l *LbsService) ListGeofences(ctx context.Context, userId string, projectId string, fenceIds []string, objects []stirng) ([]*Geofence, error) {
	logr.Debugf("ListGeofence (%s)", fenceIds)

	fences, err := l.Proxy.ListGeofence(
		auth.NewPrincipal(userId, projectId),
		fenceIds, objects)
	if err != nil {
		logr.WithError(err).Errorf("list geofence failed")
		return nil, errors.Wrap("list geofence failed", err)
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

func (l *LbsService) AddMonitoredObject(ctx context.Context, userId string, projectId string, fenceId string, objects []stirng) error {
	logr.Debugf("AddMonitoredObject (%s)", fenceId)

	err := l.Proxy.AddMonitoredObject(
		auth.NewPrincipal(userId, projectId),
		fenceId, objects)
	if err != nil {
		logr.WithError(err).Errorf("add monitored object failed")
		return errors.Wrap("ladd monitored object failed", err)
	}
	return nil
}

func (l *LbsService) RemoveMonitoredObject(ctx context.Context, userId string, projectId string, fenceId string, objects []stirng) error {
	logr.Debugf("RemoveMonitoredObject(%s)", fenceId)

	err := l.Proxy.RemoveMonitoredObject(
		auth.NewPrincipal(userId, projectId),
		fenceId, objects)
	if err != nil {
		logr.WithError(err).Errorf("remove monitored object failed")
		return errors.Wrap("remove monitored object failed", err)
	}
	return nil
}

func (l *LbsService) ListMonitoredObjects(ctx context.Context, userId string, projectId string, fenceId string, pageIndex int32, pageSize int32) (int32, []string, error) {
	logr.Debugf("ListMonitoredObjects(%s)", fenceId)

	total, objects := l.Proxy.ListMonitoredObjects(
		auth.NewPrincipal(userId, projectId),
		fenceId, int(pageIndex), int(pageSize))

	return int32(total), objects, nil
}

func (l *LbsService) ListCollections(ctx context.Context, userId string) ([]string, error) {
	logr.Debugln("ListCollections")

	//	productList := getProductList(products)
	//	logr.Debugln("productList:", productList)

	//	return productList, nil
	return nil, nil
}

func (l *LbsService) GetFenceIds(ctx context.Context, userId string, projectId string, fenceIds []string) ([]string], error) {
	logr.Debugf("GetFenceIds (%s)", in.String())
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

func (l *LbsService) QueryStatus(ctx context.Context, userId string, projectId string, monitoredPerson string, fendeIds []string) (*QueryStatus, error) {
	logr.Debugf("QueryStatus(%s)", fenceIds)

	fenceStatus, err := l.Proxy.QueryStatus(
		auth.NewPrincipal(userId, projectId),
		monitoredPerson, fenceIds)
	if err != nil {
		logr.Errorln("QueryStatus failed:", err)
		return nil, errors.Wrap("not found", err)
	}

	rsp := getMonitoredStatus(fenceStatus)

	return rsp, nil
}

func getMonitoredStatus(fenceStatus lbp.BaiduQueryStatusResponse) *QueryStatus {
	rsp := &QueryStatus{
		Status:  int32(fenceStatus.Status),
		Message: fenceStatus.Message,
		Size:    int32(fenceStatus.Size),
	}
	for _, mpVal := range fenceStatus.MonitoredStatuses {
		monitoredStatus := &pb.MonitoredStatus{
			FenceId:         int32(mpVal.FenceId),
			MonitoredStatus: mpVal.MonitoredStatus,
		}
		rsp.MonitoredStatuses = append(rsp.MonitoredStatuses, monitoredStatus)
	}
	logr.Debugln("get monitered status:", rsp)
	return rsp
}

func (l *LbsService) GetHistoryAlarms(ctx context.Context, userId string, projectId string, monitoredPerson string, fendeIds []string) (*HistoryAlarms, error) {
	logr.Debugf("GetHistoryAlarms (%s)", fendeIds)

	alarmPoint, err := l.Proxy.GetHistoryAlarms(
		auth.NewPrincipal(userId, projectId),
		monitoredPerson, fenceIds)
	if err != nil {
		logr.Errorln("GetHistoryAlarms failed:", err)
		return nil, errors.Wrap("not found", err)
	}

	rsp := getHistoryAlarmPoint(alarmPoint)

	return rsp, nil
}

func getHistoryAlarmPoint(alarmPoint lbp.BaiduGetHistoryAlarmsResponse) *HistoryAlarms {
	rsp := &HistoryAlarms{
		Status:  int32(alarmPoint.Status),
		Message: alarmPoint.Message,
		Size:    int32(alarmPoint.Size),
	}
	for _, haVal := range alarmPoint.Alarms {

		alarm := &Alarm{
			FenceId:         int32(haVal.FenceId),
			FenceName:       haVal.FenceName,
			MonitoredPerson: haVal.MonitoredPerson,
			Action:          haVal.Action,
			AlarmPoint: &AlarmPoint{
				Longitude:  haVal.AlarmPoint.Longitude,
				Latitude:   haVal.AlarmPoint.Latitude,
				Radius:     int32(haVal.AlarmPoint.Radius),
				CoordType:  haVal.AlarmPoint.CoordType,
				LocTime:    haVal.AlarmPoint.LocTime,
				CreateTime: haVal.AlarmPoint.CreateTime,
			},
			PrePoint: &PrePoint{
				Longitude:  haVal.AlarmPoint.Longitude,
				Latitude:   haVal.AlarmPoint.Latitude,
				Radius:     int32(haVal.AlarmPoint.Radius),
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

func (l *LbsService) BatchGetHistoryAlarms(ct context.Context, userId string, projectId string, coordTypeOutput string, endTime string, startTime string, pageIndex int32, pageSize int32) (*HistoryAlarms, error) {
	logr.Debugf("RemoveCollection ")

	historyAlarms, err := l.Proxy.BatchGetHistoryAlarms(
		auth.NewPrincipal(userId, projectId),
		coordTypeOutput, endTime, startTime, pageIndex, pageSize)
	if err != nil {
		logr.Errorln("BatchGetHistoryAlarms failed:", err)
		return nil, errors.Wrap("not found", err)
	}

	rsp := getBatchHistoryAlarmPoint(historyAlarms)

	return rsp, nil

}

func getBatchHistoryAlarmPoint(historyAlarms lbp.BaiduBatchHistoryAlarmsResp) *HistoryAlarms {
	rsp := &HistoryAlarms{
		Status:  int32(historyAlarms.Status),
		Message: historyAlarms.Message,
		Size:    int32(historyAlarms.Size),
		Total:   int32(historyAlarms.Total),
	}
	for _, haVal := range historyAlarms.Alarms {

		alarm := &Alarm{
			FenceId:         int32(haVal.FenceId),
			FenceName:       haVal.FenceName,
			MonitoredPerson: haVal.MonitoredPerson,
			Action:          haVal.Action,
			AlarmPoint: &AlarmPoint{
				Longitude:  haVal.AlarmPoint.Longitude,
				Latitude:   haVal.AlarmPoint.Latitude,
				Radius:     int32(haVal.AlarmPoint.Radius),
				CoordType:  haVal.AlarmPoint.CoordType,
				LocTime:    haVal.AlarmPoint.LocTime,
				CreateTime: haVal.AlarmPoint.CreateTime,
			},
			PrePoint: &PrePoint{
				Longitude:  haVal.AlarmPoint.Longitude,
				Latitude:   haVal.AlarmPoint.Latitude,
				Radius:     int32(haVal.AlarmPoint.Radius),
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

func (l *LbsService) GetStayPoints(ctx context.Context, userId string, projectId string, coordTypeOutput string, endTime string, startTime string, pageIndex int32, pageSize int32, fenceIds []string, entityName string) (*StayPoints, error) {
	logr.Debugf("GetStayPoints ()")

	stayPoints, err := l.Proxy.GetStayPoints(
		auth.NewPrincipal(userId, projectId),
		coordTypeOutput,endTime,startTime,pageIndex,pageSize,fenceIds,entityName)
	if err != nil {
		logr.Errorln("BatchGetHistoryAlarms failed:", err)
		return nil, errors.Wrap("not found", err)
	}

	rsp := getGrpcStayPoints(stayPoints)

	return rsp, nil
}

func getGrpcStayPoints(stayPoints lbp.BaiduGetStayPointResp) *StayPoints {
	rsp := &GetStayPoints{
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

func (l *LbsService) NotifyAlarms(ctx context.Context, userId string, projectId string, content string) error {
	logr.Debugf("NotifyAlarms(%s)", content)

	alarm, err := l.Proxy.UnmarshalAlarmNotification(
		auth.NewPrincipal(userId, projectId),
		content)
	if err != nil {
		logr.WithError(err).Errorf("unmarshal alarm failed")
		return errors.Wrap("unmarshal alarm failed", err)
	}
	config := message.NewConfigWithViper()
	producer, err := message.NewProducer(config, false)
	if err != nil {
		logr.WithError(err).Errorf("create message producer failed")
		return errors.Wrap("create message producer failed", err)
	}

	if err := producer.SendMessage(&lbp.AlarmTopic{Alarm: alarm}); err != nil {
		logr.WithError(err).Errorf("send alarm failed")
		return errors.Wrap("send alarm failed", err)
	}
	return nil
}

func (l *LbsService) GetFenceUserId(ctx context.Context, fenceId string) (string, error) {
	/*
		userId, err := l.Proxy.GetFenceUserId(nil, fenceId)
		if err != nil {
			logr.WithError(err).Errorf("add collection '%s' failed", fenceId)
			return nil, err
		}
		return userId, nil
	*/
	return nil, nil
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

func (l *LbsService) AddEntity(ctx context.Context, userId string, projectId string, entityName string, entityDesc string) error {
	logr.Debugf("AddEntity (%s)", entityName)

	err := l.Proxy.AddEntity(
		auth.NewPrincipal(userId, projectId),
		entityName, entityDesc)
	if err != nil {
		logr.WithError(err).Errorf("AddEntity failed")
		return errors.Wrap("AddEntity failed", err)
	}
	return nil
}

func (l *LbsService) DeleteEntity(ctx context.Context, userId string, projectId string, entityName string) error {
	logr.Debugf("DeleteEntity (%s)", entityName)

	err := l.Proxy.DeleteEntity(
		auth.NewPrincipal(userId, projectId),
		entityName)
	if err != nil {
		logr.WithError(err).Errorf("DeleteEntity failed")
		return errors.Wrap("DeleteEntity failed", err)
	}
	return nil
}

func (l *LbsService) UpdateEntity(ctx context.Context, userId string, projectId string, entityName string, entityDesc string) error {
	logr.Debugf("UpdateEntity (%s)", entityName)

	err := l.Proxy.UpdateEntity(
		auth.NewPrincipal(userId, projectId),
		entityName, entityDesc)
	if err != nil {
		logr.WithError(err).Errorf("UpdateEntity failed")
		return errors.Wrap("UpdateEntity failed", err)
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
type ListEntityRequest struct {
	UserId               string   `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty", bson:"user_id,omitempty"`
	ProjectId            string   `protobuf:"bytes,2,opt,name=project_id,json=projectId" json:"project_id,omitempty", bson:"project_id,omitempty"`
	CoordTypeOutput      string   `protobuf:"bytes,3,opt,name=coord_type_output,json=coordTypeOutput" json:"coord_type_output,omitempty", bson:"coord_type_output,omitempty"`
	PageIndex            int32    `protobuf:"varint,4,opt,name=page_index,json=pageIndex" json:"page_index,omitempty", bson:"page_index,omitempty"`
	PageSize             int32    `protobuf:"varint,5,opt,name=page_size,json=pageSize" json:"page_size,omitempty", bson:"page_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
type ListEntityResponse struct {
	Total                int32         `protobuf:"varint,1,opt,name=total" json:"total,omitempty", bson:"total,omitempty"`
	EntityInfo           []*EntityInfo `protobuf:"bytes,2,rep,name=entity_info,json=entityInfo" json:"entity_info,omitempty", bson:"entity_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}


func (l *LbsService) ListEntity(ctx context.Context, userId string, projectId string, coordTypeOutput string, pageIndex int32, pageSize int32) (int32, []*EntityInfo, error) {
	logr.Debugf("ListEntity (%s)", coordTypeOutput)
	/*
		entitiesInfo, err := l.Proxy.GetEntity(
			auth.NewPrincipal(userId, projectId))
		if err != nil {
			logr.WithError(err).Errorf("mongo get entities info failed")
			return nil, nil, errors.Wrap("mongo get entities info faile", err)
		}
		entitiesName := getEntitiesName(entitiesInfo)

		total, entityInfo := l.Proxy.ListEntity(userId, projectId, coordTypeOutput, pageIndex, pageSize)
		if total == -1 {
			logr.Errorf("ListEntity failed")
			return nil, nil, errors.Wrap("ListEntity failed", err)
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
	return nil, nil, nil
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

