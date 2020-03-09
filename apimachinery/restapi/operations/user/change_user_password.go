// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"

	"github.com/cloustone/pandas/models"
)

// ChangeUserPasswordHandlerFunc turns a function with the right signature into a change user password handler
type ChangeUserPasswordHandlerFunc func(ChangeUserPasswordParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn ChangeUserPasswordHandlerFunc) Handle(params ChangeUserPasswordParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// ChangeUserPasswordHandler interface for that can handle valid change user password params
type ChangeUserPasswordHandler interface {
	Handle(ChangeUserPasswordParams, *models.Principal) middleware.Responder
}

// NewChangeUserPassword creates a new http.Handler for the change user password operation
func NewChangeUserPassword(ctx *middleware.Context, handler ChangeUserPasswordHandler) *ChangeUserPassword {
	return &ChangeUserPassword{Context: ctx, Handler: handler}
}

/*ChangeUserPassword swagger:route PATCH /users/{userId} User changeUserPassword

change user's password

This can only be done by the logged in user.

*/
type ChangeUserPassword struct {
	Context *middleware.Context
	Handler ChangeUserPasswordHandler
}

func (o *ChangeUserPassword) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewChangeUserPasswordParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// ChangeUserPasswordBody change user password body
// swagger:model ChangeUserPasswordBody
type ChangeUserPasswordBody struct {

	// password
	Password string `json:"password,omitempty"`
}

// Validate validates this change user password body
func (o *ChangeUserPasswordBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ChangeUserPasswordBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ChangeUserPasswordBody) UnmarshalBinary(b []byte) error {
	var res ChangeUserPasswordBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}