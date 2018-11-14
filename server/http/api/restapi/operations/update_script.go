// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"
	validate "github.com/go-openapi/validate"
)

// UpdateScriptHandlerFunc turns a function with the right signature into a update script handler
type UpdateScriptHandlerFunc func(UpdateScriptParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateScriptHandlerFunc) Handle(params UpdateScriptParams) middleware.Responder {
	return fn(params)
}

// UpdateScriptHandler interface for that can handle valid update script params
type UpdateScriptHandler interface {
	Handle(UpdateScriptParams) middleware.Responder
}

// NewUpdateScript creates a new http.Handler for the update script operation
func NewUpdateScript(ctx *middleware.Context, handler UpdateScriptHandler) *UpdateScript {
	return &UpdateScript{Context: ctx, Handler: handler}
}

/*UpdateScript swagger:route PUT /projects/{projectId}/scripts/{scriptId} updateScript

Update Script

*/
type UpdateScript struct {
	Context *middleware.Context
	Handler UpdateScriptHandler
}

func (o *UpdateScript) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateScriptParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// UpdateScriptBody Script Common
//
// The properties that are shared amongst all versions of the Script model.
// swagger:model UpdateScriptBody
type UpdateScriptBody struct {

	// description
	// Max Length: 255
	// Min Length: 10
	Description string `json:"description,omitempty"`

	// execution
	// Required: true
	Execution *UpdateScriptParamsBodyExecution `json:"execution"`

	// name
	// Required: true
	// Max Length: 100
	// Min Length: 3
	Name *string `json:"name"`

	// persistence
	// Required: true
	Persistence *UpdateScriptParamsBodyPersistence `json:"persistence"`
}

// Validate validates this update script body
func (o *UpdateScriptBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateExecution(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validatePersistence(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateScriptBody) validateDescription(formats strfmt.Registry) error {

	if swag.IsZero(o.Description) { // not required
		return nil
	}

	if err := validate.MinLength("body"+"."+"description", "body", string(o.Description), 10); err != nil {
		return err
	}

	if err := validate.MaxLength("body"+"."+"description", "body", string(o.Description), 255); err != nil {
		return err
	}

	return nil
}

func (o *UpdateScriptBody) validateExecution(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"execution", "body", o.Execution); err != nil {
		return err
	}

	if o.Execution != nil {
		if err := o.Execution.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "execution")
			}
			return err
		}
	}

	return nil
}

func (o *UpdateScriptBody) validateName(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"name", "body", o.Name); err != nil {
		return err
	}

	if err := validate.MinLength("body"+"."+"name", "body", string(*o.Name), 3); err != nil {
		return err
	}

	if err := validate.MaxLength("body"+"."+"name", "body", string(*o.Name), 100); err != nil {
		return err
	}

	return nil
}

func (o *UpdateScriptBody) validatePersistence(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"persistence", "body", o.Persistence); err != nil {
		return err
	}

	if o.Persistence != nil {
		if err := o.Persistence.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "persistence")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *UpdateScriptBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateScriptBody) UnmarshalBinary(b []byte) error {
	var res UpdateScriptBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// UpdateScriptOKBody Entity
//
// Represents a database entity
// swagger:model UpdateScriptOKBody
type UpdateScriptOKBody struct {

	// id
	// Required: true
	ID *string `json:"id"`

	// rev
	// Required: true
	Rev *string `json:"rev"`

	// created at
	// Required: true
	CreatedAt *string `json:"createdAt"`

	// updated at
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *UpdateScriptOKBody) UnmarshalJSON(raw []byte) error {
	// UpdateScriptOKBodyAO0
	var dataUpdateScriptOKBodyAO0 struct {
		ID *string `json:"id"`

		Rev *string `json:"rev"`
	}
	if err := swag.ReadJSON(raw, &dataUpdateScriptOKBodyAO0); err != nil {
		return err
	}

	o.ID = dataUpdateScriptOKBodyAO0.ID

	o.Rev = dataUpdateScriptOKBodyAO0.Rev

	// UpdateScriptOKBodyAO1
	var dataUpdateScriptOKBodyAO1 struct {
		CreatedAt *string `json:"createdAt"`

		UpdatedAt string `json:"updatedAt,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataUpdateScriptOKBodyAO1); err != nil {
		return err
	}

	o.CreatedAt = dataUpdateScriptOKBodyAO1.CreatedAt

	o.UpdatedAt = dataUpdateScriptOKBodyAO1.UpdatedAt

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o UpdateScriptOKBody) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	var dataUpdateScriptOKBodyAO0 struct {
		ID *string `json:"id"`

		Rev *string `json:"rev"`
	}

	dataUpdateScriptOKBodyAO0.ID = o.ID

	dataUpdateScriptOKBodyAO0.Rev = o.Rev

	jsonDataUpdateScriptOKBodyAO0, errUpdateScriptOKBodyAO0 := swag.WriteJSON(dataUpdateScriptOKBodyAO0)
	if errUpdateScriptOKBodyAO0 != nil {
		return nil, errUpdateScriptOKBodyAO0
	}
	_parts = append(_parts, jsonDataUpdateScriptOKBodyAO0)

	var dataUpdateScriptOKBodyAO1 struct {
		CreatedAt *string `json:"createdAt"`

		UpdatedAt string `json:"updatedAt,omitempty"`
	}

	dataUpdateScriptOKBodyAO1.CreatedAt = o.CreatedAt

	dataUpdateScriptOKBodyAO1.UpdatedAt = o.UpdatedAt

	jsonDataUpdateScriptOKBodyAO1, errUpdateScriptOKBodyAO1 := swag.WriteJSON(dataUpdateScriptOKBodyAO1)
	if errUpdateScriptOKBodyAO1 != nil {
		return nil, errUpdateScriptOKBodyAO1
	}
	_parts = append(_parts, jsonDataUpdateScriptOKBodyAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this update script o k body
func (o *UpdateScriptOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateRev(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateScriptOKBody) validateID(formats strfmt.Registry) error {

	if err := validate.Required("updateScriptOK"+"."+"id", "body", o.ID); err != nil {
		return err
	}

	return nil
}

func (o *UpdateScriptOKBody) validateRev(formats strfmt.Registry) error {

	if err := validate.Required("updateScriptOK"+"."+"rev", "body", o.Rev); err != nil {
		return err
	}

	return nil
}

func (o *UpdateScriptOKBody) validateCreatedAt(formats strfmt.Registry) error {

	if err := validate.Required("updateScriptOK"+"."+"createdAt", "body", o.CreatedAt); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *UpdateScriptOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateScriptOKBody) UnmarshalBinary(b []byte) error {
	var res UpdateScriptOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// UpdateScriptParamsBodyExecution Script Execution Settings
//
// Represents script execution settings like query and params
// swagger:model UpdateScriptParamsBodyExecution
type UpdateScriptParamsBodyExecution struct {

	// params
	Params map[string]interface{} `json:"params,omitempty"`

	// query
	// Required: true
	// Min Length: 8
	Query *string `json:"query"`
}

// Validate validates this update script params body execution
func (o *UpdateScriptParamsBodyExecution) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateQuery(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateScriptParamsBodyExecution) validateQuery(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"execution"+"."+"query", "body", o.Query); err != nil {
		return err
	}

	if err := validate.MinLength("body"+"."+"execution"+"."+"query", "body", string(*o.Query), 8); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *UpdateScriptParamsBodyExecution) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateScriptParamsBodyExecution) UnmarshalBinary(b []byte) error {
	var res UpdateScriptParamsBodyExecution
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// UpdateScriptParamsBodyPersistence Script Persistence
// swagger:model UpdateScriptParamsBodyPersistence
type UpdateScriptParamsBodyPersistence struct {

	// local
	Local string `json:"local,omitempty"`

	// remote
	Remote []string `json:"remote"`
}

// Validate validates this update script params body persistence
func (o *UpdateScriptParamsBodyPersistence) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *UpdateScriptParamsBodyPersistence) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateScriptParamsBodyPersistence) UnmarshalBinary(b []byte) error {
	var res UpdateScriptParamsBodyPersistence
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}