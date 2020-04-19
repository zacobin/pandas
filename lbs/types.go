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
	Name             string    `json:"name"`
	MonitoredObjects []string  `json:"monitored_objects"`
	Longitude        float64   `json:"longitude"`
	Latitude         float64   `json:"latitude"`
	Radius           float64   `json:"radius"`
	CoordType        CoordType `json:"coord_type"`
	Denoise          int       `json:"denoise"`
	FenceID          string    `json:"fencie_id"`
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
	Name             string    `json:"name"`
	MonitoredObjects []string  `json:"monitored_objects"`
	Vertexes         string    `json:"vertexes"`
	CoordType        CoordType `json:"coord_type"`
	Denoise          int       `json:"denoise"`
	FenceID          string    `json:"fencie_id"`
}

type PrePoint struct {
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Radius     int     `json:"radius"`
	CoordType  string  `json:"coord_type"`
	LocTime    string  `json:"loctime"`
	CreateTime string  `json:"create_time"`
}

type EntityInfo struct {
	EntityName string  `json:"entity_name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
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
	Status            int               `json:"status"`
	Message           string            `json:"message"`
	Size              int               `json:"size"`
	MonitoredStatuses []MonitoredStatus `json:"monitored+statuses"`
}

type StayPoints struct {
	Status     int      `json:"status"`
	Message    string   `json:"message"`
	Total      int      `json:"total"`
	Size       int      `json:"size"`
	Distance   int      `json:"distance"`
	EndPoint   *Point   `json:"end_point"`
	StartPoint *Point   `json:"start_point"`
	Points     []*Point `json:"points"`
}

type HistoryAlarms struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Total   int      `json:"total"`
	Size    int      `json:"size"`
	Alarms  []*Alarm `json:"alarms"`
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
	EndTime         string   `json:"end_time,omitempty"`
	EntityName      string   `json:"entity_name,omitempty"`
	FenceIDs        []string `json:"fence_ids,omitempty"`
	PageIndex       int      `json:"page_index,omitempty"`
	PageSize        int      `json:"page_size,omitempty"`
	StartTime       string   `json:"start_time,omitempty"`
	CoordTypeOutput string   `json:"coord_type_output,omitempty"`
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
