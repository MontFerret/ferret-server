// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/MontFerret/ferret-server/server/http/api/models"
)

// GetScriptOKCode is the HTTP code returned for type GetScriptOK
const GetScriptOKCode int = 200

/*GetScriptOK get script o k

swagger:response getScriptOK
*/
type GetScriptOK struct {

	/*
	  In: Body
	*/
	Payload *models.ScriptOutputDetailed `json:"body,omitempty"`
}

// NewGetScriptOK creates GetScriptOK with default headers values
func NewGetScriptOK() *GetScriptOK {

	return &GetScriptOK{}
}

// WithPayload adds the payload to the get script o k response
func (o *GetScriptOK) WithPayload(payload *models.ScriptOutputDetailed) *GetScriptOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get script o k response
func (o *GetScriptOK) SetPayload(payload *models.ScriptOutputDetailed) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetScriptOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
