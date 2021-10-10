package http

import (
	"fmt"
	"net/http"

	"github.com/cloustone/pandas/mainflux"
)

var (
	_ mainflux.Response = (*removeRes)(nil)
	_ mainflux.Response = (*streamRes)(nil)
	_ mainflux.Response = (*viewStreamRes)(nil)
	_ mainflux.Response = (*streamsPageRes)(nil)
	_ mainflux.Response = (*ruleRes)(nil)
	_ mainflux.Response = (*rulesPageRes)(nil)
	_ mainflux.Response = (*ruleControlRes)(nil)
)

type removeRes struct{}

func (res removeRes) Code() int {
	return http.StatusNoContent
}

func (res removeRes) Headers() map[string]string {
	return map[string]string{}
}

func (res removeRes) Empty() bool {
	return true
}

type streamRes struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name,omitempty"`
	Json     string                 `json:"json"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	created  bool
}

func (res streamRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res streamRes) Headers() map[string]string {
	if res.created {
		return map[string]string{
			"Location":           fmt.Sprintf("/streams/%s", res.ID),
			"Warning-Deprecated": "This endpoint will be depreciated in v1.0.0. It will be replaced with the bulk endpoint currently found at /streams/bulk.",
		}
	}

	return map[string]string{}
}

func (res streamRes) Empty() bool {
	return true
}

type streamsRes struct {
	Streams []streamRes `json:"streams"`
	created bool
}

func (res streamsRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res streamsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res streamsRes) Empty() bool {
	return false
}

type viewStreamRes struct {
	ID       string                 `json:"id"`
	Owner    string                 `json:"-"`
	Name     string                 `json:"name,omitempty"`
	Json     string                 `json:"json"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (res viewStreamRes) Code() int {
	return http.StatusOK
}

func (res viewStreamRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewStreamRes) Empty() bool {
	return false
}

type streamsPageRes struct {
	pageRes
	Streams []viewStreamRes `json:"streams"`
}

func (res streamsPageRes) Code() int {
	return http.StatusOK
}

func (res streamsPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res streamsPageRes) Empty() bool {
	return false
}

type ruleRes struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name,omitempty"`
	SQL      string                 `json:"sql,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	created  bool
}

func (res ruleRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res ruleRes) Headers() map[string]string {
	if res.created {
		return map[string]string{
			"Location":           fmt.Sprintf("/rules/%s", res.ID),
			"Warning-Deprecated": "This endpoint will be depreciated in v1.0.0. It will be replaced with the bulk endpoint currently found at /rules/bulk.",
		}
	}

	return map[string]string{}
}

func (res ruleRes) Empty() bool {
	return true
}

type rulesRes struct {
	Rules   []ruleRes `json:"rules"`
	created bool
}

func (res rulesRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res rulesRes) Headers() map[string]string {
	return map[string]string{}
}

func (res rulesRes) Empty() bool {
	return false
}

type rulesPageRes struct {
	pageRes
	Rules []ruleRes `json:"rules"`
}

func (res rulesPageRes) Code() int {
	return http.StatusOK
}

func (res rulesPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res rulesPageRes) Empty() bool {
	return false
}

type connectionRes struct{}

func (res connectionRes) Code() int {
	return http.StatusOK
}

func (res connectionRes) Headers() map[string]string {
	return map[string]string{
		"Warning-Deprecated": "This endpoint will be depreciated in v1.0.0. It will be replaced with the bulk endpoint found at /connect.",
	}
}

func (res connectionRes) Empty() bool {
	return true
}

type createConnectionsRes struct{}

func (res createConnectionsRes) Code() int {
	return http.StatusOK
}

func (res createConnectionsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res createConnectionsRes) Empty() bool {
	return true
}

type ruleControlRes struct{}

func (res ruleControlRes) Code() int {
	return http.StatusNoContent
}

func (res ruleControlRes) Headers() map[string]string {
	return map[string]string{}
}

func (res ruleControlRes) Empty() bool {
	return true
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}
