// SPDX-License-Identifier: Apache-2.0
package alerts

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type AlertMethod string

const (
	AlertMethodEmail   Method = "email"
	AlertMethodMessage Method = "message"
	AlertMethodTel     Method = "tel"
)

// Metadata stores arbitrary variable data
type Metadata map[string]interface{}

// Alert
type Alert struct {
	ID        string          `json:"id,omitempty"`
	Owner     string          `json:"owner"`
	Name      string          `json:"name,omitempty"`
	Group     string          `json:"group,omitempty"`
	Time      string          `json:"time,omitempty"`
	On        bool            `json:"on,omitempty"`
	Methods   []AlertMethod   `json:"method,omitempty"`
	Users     []AlertUser     `json:"users,omitempty""`
	Metadata  Metadata        `json:"metadta, omitempty"`
	Extension *AlertExtension `json:"alert_extension"`
}

type AlertUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Roles     []string  `json:"roles"`
	Status    string    `json:"status"`
}

type AlertExtension struct {
	Labels       map[string]string
	Annotations  map[string]string
	StartAt      time.Time
	EndAt        time.Time
	GeneratorUrl string
	State        string
	SilienceBy   []string
	InhibitBy    []string
}

type AlertRule struct {
	Alert       string
	Expr        string
	For         string
	Labels      map[string]string
	Annotations map[string]string
}

type AlertGroup struct {
	Name       string
	AlertRules []AlertRule
}

type AlertRules struct {
	AlertGroups []AlertGroup
}

type AlertRuleExtension struct {
	Alert       string
	Expr        string
	For         string
	Labels      map[string]string
	Group       string
	Annotations map[string]string
	Labels      map[string]string
}

type Point struct {
	Longitude float64              `json:"longitude,omitempty `
	Latitude  float64              `json:"latitude,omitempty"`
	Radius    int32                `json:"radius,omitempty"`
	CoordType string               `json:"coord_type,omitempty"`
	LocTime   *timestamp.Timestamp `json:"loc_time,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"create_at,omitempty"`
}

type Alarm struct {
	FenceID          int32  `json:"id,omitempty"`
	FenceName        string `json:"fence_name,omitempty"`
	MonitoredObjects string `json:"monitored_objects,omitempty"`
	UserID           string `json:"user_id,omitempty"`
	Action           string `json:"action,omitempty" bson:"action,omitempty"`
	AlarmPoint       *Point `json:"alarm_point,omitempty"`
	PrePoint         *Point `json:"pre_point,omitempty"`
}
