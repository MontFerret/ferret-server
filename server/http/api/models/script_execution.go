// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ScriptExecution Script Execution Settings
//
// Represents script execution settings like query and params
// swagger:model script-execution
type ScriptExecution struct {

	// params
	Params map[string]Any `json:"params,omitempty"`

	// query
	// Required: true
	// Min Length: 8
	Query *string `json:"query"`
}

// Validate validates this script execution
func (m *ScriptExecution) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateQuery(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ScriptExecution) validateQuery(formats strfmt.Registry) error {

	if err := validate.Required("query", "body", m.Query); err != nil {
		return err
	}

	if err := validate.MinLength("query", "body", string(*m.Query), 8); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ScriptExecution) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ScriptExecution) UnmarshalBinary(b []byte) error {
	var res ScriptExecution
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
