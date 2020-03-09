// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewChangeUserPasswordParams creates a new ChangeUserPasswordParams object
// no default values defined in spec.
func NewChangeUserPasswordParams() ChangeUserPasswordParams {

	return ChangeUserPasswordParams{}
}

// ChangeUserPasswordParams contains all the bound params for the change user password operation
// typically these are obtained from a http.Request
//
// swagger:parameters ChangeUserPassword
type ChangeUserPasswordParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Updated user password
	  Required: true
	  In: body
	*/
	User ChangeUserPasswordBody
	/*name that need to be updated
	  Required: true
	  In: path
	*/
	UserID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewChangeUserPasswordParams() beforehand.
func (o *ChangeUserPasswordParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body ChangeUserPasswordBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("user", "body"))
			} else {
				res = append(res, errors.NewParseError("user", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.User = body
			}
		}
	} else {
		res = append(res, errors.Required("user", "body"))
	}
	rUserID, rhkUserID, _ := route.Params.GetOK("userId")
	if err := o.bindUserID(rUserID, rhkUserID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindUserID binds and validates parameter UserID from path.
func (o *ChangeUserPasswordParams) bindUserID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.UserID = raw

	return nil
}