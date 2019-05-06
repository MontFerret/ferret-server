// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewFindExecutionsParams creates a new FindExecutionsParams object
// with the default values initialized.
func NewFindExecutionsParams() FindExecutionsParams {

	var (
		// initialize parameters with default values

		countDefault = int32(10)
	)

	return FindExecutionsParams{
		Count: &countDefault,
	}
}

// FindExecutionsParams contains all the bound params for the find executions operation
// typically these are obtained from a http.Request
//
// swagger:parameters findExecutions
type FindExecutionsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Script execution cause
	  In: query
	*/
	Cause *string
	/*Count of items to return
	  Maximum: 100
	  Minimum: 1
	  In: query
	  Default: 10
	*/
	Count *int32
	/*Pagination cursor
	  In: query
	*/
	Cursor *int64
	/*
	  Required: true
	  Pattern: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}
	  In: path
	*/
	ProjectID string
	/*
	  In: query
	*/
	Status *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewFindExecutionsParams() beforehand.
func (o *FindExecutionsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qCause, qhkCause, _ := qs.GetOK("cause")
	if err := o.bindCause(qCause, qhkCause, route.Formats); err != nil {
		res = append(res, err)
	}

	qCount, qhkCount, _ := qs.GetOK("count")
	if err := o.bindCount(qCount, qhkCount, route.Formats); err != nil {
		res = append(res, err)
	}

	qCursor, qhkCursor, _ := qs.GetOK("cursor")
	if err := o.bindCursor(qCursor, qhkCursor, route.Formats); err != nil {
		res = append(res, err)
	}

	rProjectID, rhkProjectID, _ := route.Params.GetOK("projectID")
	if err := o.bindProjectID(rProjectID, rhkProjectID, route.Formats); err != nil {
		res = append(res, err)
	}

	qStatus, qhkStatus, _ := qs.GetOK("status")
	if err := o.bindStatus(qStatus, qhkStatus, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCause binds and validates parameter Cause from query.
func (o *FindExecutionsParams) bindCause(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Cause = &raw

	if err := o.validateCause(formats); err != nil {
		return err
	}

	return nil
}

// validateCause carries on validations for parameter Cause
func (o *FindExecutionsParams) validateCause(formats strfmt.Registry) error {

	if err := validate.Enum("cause", "query", *o.Cause, []interface{}{"manual", "schedule", "hook", "unknown"}); err != nil {
		return err
	}

	return nil
}

// bindCount binds and validates parameter Count from query.
func (o *FindExecutionsParams) bindCount(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewFindExecutionsParams()
		return nil
	}

	value, err := swag.ConvertInt32(raw)
	if err != nil {
		return errors.InvalidType("count", "query", "int32", raw)
	}
	o.Count = &value

	if err := o.validateCount(formats); err != nil {
		return err
	}

	return nil
}

// validateCount carries on validations for parameter Count
func (o *FindExecutionsParams) validateCount(formats strfmt.Registry) error {

	if err := validate.MinimumInt("count", "query", int64(*o.Count), 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("count", "query", int64(*o.Count), 100, false); err != nil {
		return err
	}

	return nil
}

// bindCursor binds and validates parameter Cursor from query.
func (o *FindExecutionsParams) bindCursor(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("cursor", "query", "int64", raw)
	}
	o.Cursor = &value

	return nil
}

// bindProjectID binds and validates parameter ProjectID from path.
func (o *FindExecutionsParams) bindProjectID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ProjectID = raw

	if err := o.validateProjectID(formats); err != nil {
		return err
	}

	return nil
}

// validateProjectID carries on validations for parameter ProjectID
func (o *FindExecutionsParams) validateProjectID(formats strfmt.Registry) error {

	if err := validate.Pattern("projectID", "path", o.ProjectID, `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`); err != nil {
		return err
	}

	return nil
}

// bindStatus binds and validates parameter Status from query.
func (o *FindExecutionsParams) bindStatus(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Status = &raw

	if err := o.validateStatus(formats); err != nil {
		return err
	}

	return nil
}

// validateStatus carries on validations for parameter Status
func (o *FindExecutionsParams) validateStatus(formats strfmt.Registry) error {

	if err := validate.Enum("status", "query", *o.Status, []interface{}{"unknown", "queued", "running", "completed", "cancelled", "errored"}); err != nil {
		return err
	}

	return nil
}
