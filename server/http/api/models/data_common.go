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

// DataCommon Data Common
//
// The properties that are shared amongst all versions of the Data model.
// swagger:model data-common
type DataCommon struct {

	// job Id
	// Required: true
	JobID *string `json:"jobId"`

	// script Id
	// Required: true
	ScriptID *string `json:"scriptId"`

	// script rev
	// Required: true
	ScriptRev *string `json:"scriptRev"`

	// value
	// Required: true
	Value Any `json:"value"`
}

// Validate validates this data common
func (m *DataCommon) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateJobID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateScriptID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateScriptRev(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataCommon) validateJobID(formats strfmt.Registry) error {

	if err := validate.Required("jobId", "body", m.JobID); err != nil {
		return err
	}

	return nil
}

func (m *DataCommon) validateScriptID(formats strfmt.Registry) error {

	if err := validate.Required("scriptId", "body", m.ScriptID); err != nil {
		return err
	}

	return nil
}

func (m *DataCommon) validateScriptRev(formats strfmt.Registry) error {

	if err := validate.Required("scriptRev", "body", m.ScriptRev); err != nil {
		return err
	}

	return nil
}

func (m *DataCommon) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("value", "body", m.Value); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DataCommon) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataCommon) UnmarshalBinary(b []byte) error {
	var res DataCommon
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
