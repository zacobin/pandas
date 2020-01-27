// Code generated by go-swagger; DO NOT EDIT.

package dashboard

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/cloustone/pandas/models"
)

// GetDashboardOKCode is the HTTP code returned for type GetDashboardOK
const GetDashboardOKCode int = 200

/*GetDashboardOK Successful operation

swagger:response getDashboardOK
*/
type GetDashboardOK struct {

	/*
	  In: Body
	*/
	Payload *models.Dashboard `json:"body,omitempty"`
}

// NewGetDashboardOK creates GetDashboardOK with default headers values
func NewGetDashboardOK() *GetDashboardOK {

	return &GetDashboardOK{}
}

// WithPayload adds the payload to the get dashboard o k response
func (o *GetDashboardOK) WithPayload(payload *models.Dashboard) *GetDashboardOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get dashboard o k response
func (o *GetDashboardOK) SetPayload(payload *models.Dashboard) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDashboardOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetDashboardBadRequestCode is the HTTP code returned for type GetDashboardBadRequest
const GetDashboardBadRequestCode int = 400

/*GetDashboardBadRequest Bad request

swagger:response getDashboardBadRequest
*/
type GetDashboardBadRequest struct {
}

// NewGetDashboardBadRequest creates GetDashboardBadRequest with default headers values
func NewGetDashboardBadRequest() *GetDashboardBadRequest {

	return &GetDashboardBadRequest{}
}

// WriteResponse to the client
func (o *GetDashboardBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// GetDashboardInternalServerErrorCode is the HTTP code returned for type GetDashboardInternalServerError
const GetDashboardInternalServerErrorCode int = 500

/*GetDashboardInternalServerError Server internal error

swagger:response getDashboardInternalServerError
*/
type GetDashboardInternalServerError struct {
}

// NewGetDashboardInternalServerError creates GetDashboardInternalServerError with default headers values
func NewGetDashboardInternalServerError() *GetDashboardInternalServerError {

	return &GetDashboardInternalServerError{}
}

// WriteResponse to the client
func (o *GetDashboardInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}