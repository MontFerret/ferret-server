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

// ExecutionInput Execution Input
//
// The properties that are allowed when creating or updating a Execution.
// swagger:model execution-input
type ExecutionInput struct {

	// params
	Params map[string]Any `json:"params,omitempty"`

	// script ID
	// Required: true
	// Pattern: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}
	ScriptID *string `json:"scriptID"`
}

// Validate validates this execution input
func (m *ExecutionInput) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateScriptID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ExecutionInput) validateScriptID(formats strfmt.Registry) error {

	if err := validate.Required("scriptID", "body", m.ScriptID); err != nil {
		return err
	}

	if err := validate.Pattern("scriptID", "body", string(*m.ScriptID), `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ExecutionInput) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ExecutionInput) UnmarshalBinary(b []byte) error {
	var res ExecutionInput
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
