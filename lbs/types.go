package lbs

type CoordType string

const (
	CoordType_WGS84  CoordType = "wgs84"
	CoordType_GCJ02  CoordType = "gcj02"
	CoordType_BD09LL CoordType = "bd09ll"
)

type TrackPoint struct {
	EntityName string    `json:"entity_name"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	CoordType  CoordType `json:"coord_type_input"`
	Time       string    `json:"loc_time"`
}

// Geofence
type CircleGeofence struct {
	Name             string
	MonitoredObjects []string
	Longitude        float64
	Latitude         float64
	Radius           float64
	CoordType        CoordType
	Denoise          int
	FenceID          string
}

type Vertexe struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Geofence struct {
	FenceID         string    `json:"fence_id"`
	FenceName       string    `json:"fence_name"`
	MonitoredObject string    `json:"monitored_person"`
	Shape           string    `json:"shape"`
	Longitude       float64   `json:"longitude"`
	Latitude        float64   `json:"latitude"`
	Radius          float64   `json:"radius"`
	CoordType       CoordType `json:"coord_type"`
	Denoise         int       `json:"denoise"`
	CreateTime      string    `json:"create_time"`
	UpdateTime      string    `json:"modify_time"`
	Vertexes        []Vertexe `json:"vertexes"`
}

type PolyGeofence struct {
	Name             string
	MonitoredObjects []string
	Vertexes         string
	CoordType        CoordType
	Denoise          int
	FenceID          string
}

type PrePoint struct {
	Longitude  float64
	Latitude   float64
	Radius     int
	CoordType  string
	LocTime    string
	CreateTime string
}

type EntityInfo struct {
	EntityName string
	Latitude   float64
	Longitude  float64
}

type AlarmPoint struct {
	Longitude  float64   `json:"longitude"`
	Latitude   float64   `json:"latitude"`
	Radius     int       `json:"radius"`
	CoordType  CoordType `json:"coord_type"`
	LocTime    string    `json:"loc_time"`
	CreateTime string    `json:"create_time"`
}

type AlarmPointInfo struct {
	Longitude  float64   `json:"longitude"`
	Latitude   float64   `json:"latitude"`
	Radius     int       `json:"radius"`
	CoordType  CoordType `json:"coord_type"`
	LocTime    int       `json:"loc_time"`
	CreateTime int       `json:"create_time"`
}

type Alarm struct {
	FenceID          string     `json:"fence_id,noempty"`
	FenceName        string     `json:"fence_name,noempty"`
	MonitoredObjects []string   `json:"monitored_objexts"`
	Action           string     `json:"action"`
	AlarmPoint       AlarmPoint `json:"alarm_point"`
	PrePoint         AlarmPoint `json:"pre_point"`
}

type AlarmNotification struct {
	Type      int      `json:"type"`
	ServiceID string   `json:"service_id"`
	Alarms    []*Alarm `json:"alarms"`
}

type Point struct {
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	CoordType CoordType `json:"coord_type"`
	LocTime   string    `json:"loc_time"`
}

type AlarmInfos struct {
	Type      int         `json:"type"`
	ServiceId int         `json:"service_id"`
	Alarms    []AlarmInfo `json:"content"`
}

type AlarmInfo struct {
	FenceID          int            `json:"fence_id,noempty"`
	FenceName        string         `json:"fence_name,noempty"`
	MonitoredObjects string         `json:"monitored_person"`
	Action           string         `json:"action"`
	AlarmPoint       AlarmPointInfo `json:"alarm_point"`
	PrePoint         AlarmPointInfo `json:"pre_point"`
	UserID           string
}

type QueryStatus struct {
	Status            int
	Message           string
	Size              int
	MonitoredStatuses []MonitoredStatus
}

type StayPoints struct {
	Status     int
	Message    string
	Total      int
	Size       int
	Distance   int
	EndPoint   *Point
	StartPoint *Point
	Points     []*Point
}

type HistoryAlarms struct {
	Status  int
	Message string
	Total   int
	Size    int
	Alarms  []*Alarm
}

// Location is reported by iot terminal to lbs service
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Time      string  `json:"loc_time"`
}

type MonitoredStatus struct {
	FenceID         int    `json:"fence_id"`
	MonitoredStatus string `json:"monitored_status"`
}

type QueryStatusResponse struct {
	Status            int               `json:"status"`
	Message           string            `json:"message"`
	Size              int               `json:"size"`
	MonitoredStatuses []MonitoredStatus `json:"monitored_statuses"`
}

type GetStayPointResp struct {
	Status     int     `json:"status"`
	Message    string  `json:"message"`
	Size       int     `json:"size"`
	Total      int     `json:"total"`
	StartPoint Point   `json:"start_point"`
	EndPoint   Point   `json:"end_point"`
	Points     []Point `json:"points"`
}
type GetStayPointsRequest struct {
	EndTime         string   `protobuf:"bytes,3,opt,name=end_time,json=endTime" json:"end_time,omitempty", bson:"end_time,omitempty"`
	EntityName      string   `protobuf:"bytes,4,opt,name=entity_name,json=entityName" json:"entity_name,omitempty", bson:"entity_name,omitempty"`
	FenceIDs        []string `protobuf:"bytes,5,rep,name=fence_ids,json=fenceIDs" json:"fence_ids,omitempty", bson:"fence_ids,omitempty"`
	PageIndex       int      `protobuf:"varint,6,opt,name=page_index,json=pageIndex" json:"page_index,omitempty", bson:"page_index,omitempty"`
	PageSize        int      `protobuf:"varint,7,opt,name=page_size,json=pageSize" json:"page_size,omitempty", bson:"page_size,omitempty"`
	StartTime       string   `protobuf:"bytes,8,opt,name=start_time,json=startTime" json:"start_time,omitempty", bson:"start_time,omitempty"`
	CoordTypeOutput string   `protobuf:"bytes,9,opt,name=coord_type_output,json=coordTypeOutput" json:"coord_type_output,omitempty", bson:"coord_type_output,omitempty"`
}

type HistoryAlarmPoint struct {
	Longitude  float64   `json:"longitude"`
	Latitude   float64   `json:"latitude"`
	Radius     int       `json:"radius"`
	CoordType  CoordType `json:"coord_type"`
	LocTime    string    `json:"loc_time"`
	CreateTime string    `json:"create_time"`
}

type HistoryPrePoint struct {
	Longitude  float64   `json:"longitude"`
	Latitude   float64   `json:"latitude"`
	Radius     int       `json:"radius"`
	CoordType  CoordType `json:"coord_type"`
	LocTime    string    `json:"loc_time"`
	CreateTime string    `json:"create_time"`
}

type AlarmHistory struct {
	FenceID          string            `json:"fence_id"`
	FenceName        string            `json:"fence_name"`
	MonitoredObjects []string          `json:"monitored_person"`
	Action           string            `json:"action"`
	AlarmPoint       HistoryAlarmPoint `json:"alarm_point"`
	PrePoint         HistoryPrePoint   `json:"pre_point"`
}

type HistoryAlarmsResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Size    int            `json:"size"`
	Alarms  []AlarmHistory `json:"alarms"`
}

type BatchHistoryAlarmsResp struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Size    int            `json:"size"`
	Total   int            `json:"total"`
	Alarms  []AlarmHistory `json:"alarms"`
}

type LastLocationStruct struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Entity struct {
	EntityName   string             `json:"entity_name"`
	EntityDesc   string             `json:"entity_desc"`
	LastLocation LastLocationStruct `json:"latest_location"`
}

type ListEntityStruct struct {
	Status   int      `json:"status"`
	Message  string   `json:"message"`
	Total    int      `json:"total"`
	Entities []Entity `json:"entities"`
}

type GetHistoryAlarmsRequest struct {
	UserID          string   `json:"user_id,omitempty"`
	CollectionID    string   `json:"collection_id,omitempty""`
	MonitoredPerson string   `json:"monitored_person,omitempty""`
	FenceIDs        []string `json:"fence_ids,omitempty"`
}

type BatchGetHistoryAlarmsRequest struct {
	UserID          string `json:"user_id,omitempty"`
	CollectionID    string `json:"collection_id,omitempty"`
	CoordTypeOutput string `json:"coord_type_output,omitempty"`
	EndTime         string `json:"end_time,omitempty"`
	StartTime       string `json:"start_time,omitempty"`
	PageIndex       int    `json:"page_index,omitempty"`
	PageSize        int    `json:"page_size,omitempty"`
}

const (
	AlarmMessageTopic = "pandas-alarm"
)

type AlarmTopic struct {
	TopicName string
	Alarm     *AlarmNotification `json:"alarm"`
}

func (p *AlarmTopic) Topic() string        { return AlarmMessageTopic }
func (p *AlarmTopic) SetTopic(name string) {}
func (p *AlarmTopic) Serialize(opt SerializeOption) ([]byte, error) {
	return Serialize(p, opt)
}
func (p *AlarmTopic) Deserialize(buf []byte, opt SerializeOption) error {
	return Deserialize(buf, opt, p)
}
