// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"encoding/json"
	"net/http"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"
	validate "github.com/go-openapi/validate"
)

// FindExecutionsHandlerFunc turns a function with the right signature into a find executions handler
type FindExecutionsHandlerFunc func(FindExecutionsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn FindExecutionsHandlerFunc) Handle(params FindExecutionsParams) middleware.Responder {
	return fn(params)
}

// FindExecutionsHandler interface for that can handle valid find executions params
type FindExecutionsHandler interface {
	Handle(FindExecutionsParams) middleware.Responder
}

// NewFindExecutions creates a new http.Handler for the find executions operation
func NewFindExecutions(ctx *middleware.Context, handler FindExecutionsHandler) *FindExecutions {
	return &FindExecutions{Context: ctx, Handler: handler}
}

/*FindExecutions swagger:route GET /projects/{projectID}/execution findExecutions

List Execution

*/
type FindExecutions struct {
	Context *middleware.Context
	Handler FindExecutionsHandler
}

func (o *FindExecutions) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFindExecutionsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// FindExecutionsOKBodyItems0 Execution Output
//
// The properties that are included when fetching a list of Executions.
// swagger:model FindExecutionsOKBodyItems0
type FindExecutionsOKBodyItems0 struct {
	FindExecutionsOKBodyItems0AllOf0

	// Execution Cause
	//
	// Execution cause
	// Required: true
	// Enum: [unknown manual schedule hook]
	Cause *string `json:"cause"`

	// job id
	// Required: true
	// Pattern: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}
	JobID *string `json:"job_id"`

	// script id
	// Required: true
	// Pattern: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}
	ScriptID *string `json:"script_id"`

	// script rev
	// Required: true
	ScriptRev *string `json:"script_rev"`

	// Execution Status
	//
	// Execution stats
	// Required: true
	// Enum: [unknown queued running completed cancelled errored]
	Status *string `json:"status"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *FindExecutionsOKBodyItems0) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 FindExecutionsOKBodyItems0AllOf0
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	o.FindExecutionsOKBodyItems0AllOf0 = aO0

	// AO1
	var dataAO1 struct {
		Cause *string `json:"cause"`

		JobID *string `json:"job_id"`

		ScriptID *string `json:"script_id"`

		ScriptRev *string `json:"script_rev"`

		Status *string `json:"status"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	o.Cause = dataAO1.Cause

	o.JobID = dataAO1.JobID

	o.ScriptID = dataAO1.ScriptID

	o.ScriptRev = dataAO1.ScriptRev

	o.Status = dataAO1.Status

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o FindExecutionsOKBodyItems0) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(o.FindExecutionsOKBodyItems0AllOf0)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	var dataAO1 struct {
		Cause *string `json:"cause"`

		JobID *string `json:"job_id"`

		ScriptID *string `json:"script_id"`

		ScriptRev *string `json:"script_rev"`

		Status *string `json:"status"`
	}

	dataAO1.Cause = o.Cause

	dataAO1.JobID = o.JobID

	dataAO1.ScriptID = o.ScriptID

	dataAO1.ScriptRev = o.ScriptRev

	dataAO1.Status = o.Status

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this find executions o k body items0
func (o *FindExecutionsOKBodyItems0) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with FindExecutionsOKBodyItems0AllOf0

	if err := o.validateCause(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateJobID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateScriptID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateScriptRev(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var findExecutionsOKBodyItems0TypeCausePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unknown","manual","schedule","hook"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		findExecutionsOKBodyItems0TypeCausePropEnum = append(findExecutionsOKBodyItems0TypeCausePropEnum, v)
	}
}

// property enum
func (o *FindExecutionsOKBodyItems0) validateCauseEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, findExecutionsOKBodyItems0TypeCausePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *FindExecutionsOKBodyItems0) validateCause(formats strfmt.Registry) error {

	if err := validate.Required("cause", "body", o.Cause); err != nil {
		return err
	}

	// value enum
	if err := o.validateCauseEnum("cause", "body", *o.Cause); err != nil {
		return err
	}

	return nil
}

func (o *FindExecutionsOKBodyItems0) validateJobID(formats strfmt.Registry) error {

	if err := validate.Required("job_id", "body", o.JobID); err != nil {
		return err
	}

	if err := validate.Pattern("job_id", "body", string(*o.JobID), `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`); err != nil {
		return err
	}

	return nil
}

func (o *FindExecutionsOKBodyItems0) validateScriptID(formats strfmt.Registry) error {

	if err := validate.Required("script_id", "body", o.ScriptID); err != nil {
		return err
	}

	if err := validate.Pattern("script_id", "body", string(*o.ScriptID), `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`); err != nil {
		return err
	}

	return nil
}

func (o *FindExecutionsOKBodyItems0) validateScriptRev(formats strfmt.Registry) error {

	if err := validate.Required("script_rev", "body", o.ScriptRev); err != nil {
		return err
	}

	return nil
}

var findExecutionsOKBodyItems0TypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unknown","queued","running","completed","cancelled","errored"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		findExecutionsOKBodyItems0TypeStatusPropEnum = append(findExecutionsOKBodyItems0TypeStatusPropEnum, v)
	}
}

// property enum
func (o *FindExecutionsOKBodyItems0) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, findExecutionsOKBodyItems0TypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (o *FindExecutionsOKBodyItems0) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", o.Status); err != nil {
		return err
	}

	// value enum
	if err := o.validateStatusEnum("status", "body", *o.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *FindExecutionsOKBodyItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *FindExecutionsOKBodyItems0) UnmarshalBinary(b []byte) error {
	var res FindExecutionsOKBodyItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// FindExecutionsOKBodyItems0AllOf0 find executions o k body items0 all of0
// swagger:model FindExecutionsOKBodyItems0AllOf0
type FindExecutionsOKBodyItems0AllOf0 interface{}
