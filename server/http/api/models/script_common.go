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

// ScriptCommon Script Common
//
// The properties that are shared amongst all versions of the Script model.
// swagger:model script-common
type ScriptCommon struct {
	Definition

	// execution
	// Required: true
	Execution *ScriptExecution `json:"execution"`

	// persistence
	// Required: true
	Persistence *ScriptPersistence `json:"persistence"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *ScriptCommon) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 Definition
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.Definition = aO0

	// AO1
	var dataAO1 struct {
		Execution *ScriptExecution `json:"execution"`

		Persistence *ScriptPersistence `json:"persistence"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Execution = dataAO1.Execution

	m.Persistence = dataAO1.Persistence

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m ScriptCommon) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.Definition)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	var dataAO1 struct {
		Execution *ScriptExecution `json:"execution"`

		Persistence *ScriptPersistence `json:"persistence"`
	}

	dataAO1.Execution = m.Execution

	dataAO1.Persistence = m.Persistence

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this script common
func (m *ScriptCommon) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with Definition
	if err := m.Definition.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExecution(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePersistence(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ScriptCommon) validateExecution(formats strfmt.Registry) error {

	if err := validate.Required("execution", "body", m.Execution); err != nil {
		return err
	}

	if m.Execution != nil {
		if err := m.Execution.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("execution")
			}
			return err
		}
	}

	return nil
}

func (m *ScriptCommon) validatePersistence(formats strfmt.Registry) error {

	if err := validate.Required("persistence", "body", m.Persistence); err != nil {
		return err
	}

	if m.Persistence != nil {
		if err := m.Persistence.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("persistence")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ScriptCommon) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ScriptCommon) UnmarshalBinary(b []byte) error {
	var res ScriptCommon
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
