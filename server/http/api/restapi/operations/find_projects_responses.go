// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// FindProjectsOKCode is the HTTP code returned for type FindProjectsOK
const FindProjectsOKCode int = 200

/*FindProjectsOK find projects o k

swagger:response findProjectsOK
*/
type FindProjectsOK struct {

	/*
	  In: Body
	*/
	Payload *FindProjectsOKBody `json:"body,omitempty"`
}

// NewFindProjectsOK creates FindProjectsOK with default headers values
func NewFindProjectsOK() *FindProjectsOK {

	return &FindProjectsOK{}
}

// WithPayload adds the payload to the find projects o k response
func (o *FindProjectsOK) WithPayload(payload *FindProjectsOKBody) *FindProjectsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find projects o k response
func (o *FindProjectsOK) SetPayload(payload *FindProjectsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindProjectsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
