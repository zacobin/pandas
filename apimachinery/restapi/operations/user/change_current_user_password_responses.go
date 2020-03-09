// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// ChangeCurrentUserPasswordOKCode is the HTTP code returned for type ChangeCurrentUserPasswordOK
const ChangeCurrentUserPasswordOKCode int = 200

/*ChangeCurrentUserPasswordOK successful operation

swagger:response changeCurrentUserPasswordOK
*/
type ChangeCurrentUserPasswordOK struct {
}

// NewChangeCurrentUserPasswordOK creates ChangeCurrentUserPasswordOK with default headers values
func NewChangeCurrentUserPasswordOK() *ChangeCurrentUserPasswordOK {

	return &ChangeCurrentUserPasswordOK{}
}

// WriteResponse to the client
func (o *ChangeCurrentUserPasswordOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// ChangeCurrentUserPasswordBadRequestCode is the HTTP code returned for type ChangeCurrentUserPasswordBadRequest
const ChangeCurrentUserPasswordBadRequestCode int = 400

/*ChangeCurrentUserPasswordBadRequest Invalid username supplied

swagger:response changeCurrentUserPasswordBadRequest
*/
type ChangeCurrentUserPasswordBadRequest struct {
}

// NewChangeCurrentUserPasswordBadRequest creates ChangeCurrentUserPasswordBadRequest with default headers values
func NewChangeCurrentUserPasswordBadRequest() *ChangeCurrentUserPasswordBadRequest {

	return &ChangeCurrentUserPasswordBadRequest{}
}

// WriteResponse to the client
func (o *ChangeCurrentUserPasswordBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// ChangeCurrentUserPasswordNotFoundCode is the HTTP code returned for type ChangeCurrentUserPasswordNotFound
const ChangeCurrentUserPasswordNotFoundCode int = 404

/*ChangeCurrentUserPasswordNotFound User not found

swagger:response changeCurrentUserPasswordNotFound
*/
type ChangeCurrentUserPasswordNotFound struct {
}

// NewChangeCurrentUserPasswordNotFound creates ChangeCurrentUserPasswordNotFound with default headers values
func NewChangeCurrentUserPasswordNotFound() *ChangeCurrentUserPasswordNotFound {

	return &ChangeCurrentUserPasswordNotFound{}
}

// WriteResponse to the client
func (o *ChangeCurrentUserPasswordNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}