// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"
	"strconv"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"

	models "github.com/MontFerret/ferret-server/server/http/api/models"
)

// FindScriptsHandlerFunc turns a function with the right signature into a find scripts handler
type FindScriptsHandlerFunc func(FindScriptsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn FindScriptsHandlerFunc) Handle(params FindScriptsParams) middleware.Responder {
	return fn(params)
}

// FindScriptsHandler interface for that can handle valid find scripts params
type FindScriptsHandler interface {
	Handle(FindScriptsParams) middleware.Responder
}

// NewFindScripts creates a new http.Handler for the find scripts operation
func NewFindScripts(ctx *middleware.Context, handler FindScriptsHandler) *FindScripts {
	return &FindScripts{Context: ctx, Handler: handler}
}

/*FindScripts swagger:route GET /projects/{projectId}/scripts findScripts

List Script

*/
type FindScripts struct {
	Context *middleware.Context
	Handler FindScriptsHandler
}

func (o *FindScripts) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFindScriptsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// FindScriptsOKBody find scripts o k body
// swagger:model FindScriptsOKBody
type FindScriptsOKBody struct {
	models.SearchResult

	// data
	Data []*models.ScriptOutput `json:"data"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *FindScriptsOKBody) UnmarshalJSON(raw []byte) error {
	// FindScriptsOKBodyAO0
	var findScriptsOKBodyAO0 models.SearchResult
	if err := swag.ReadJSON(raw, &findScriptsOKBodyAO0); err != nil {
		return err
	}
	o.SearchResult = findScriptsOKBodyAO0

	// FindScriptsOKBodyAO1
	var dataFindScriptsOKBodyAO1 struct {
		Data []*models.ScriptOutput `json:"data"`
	}
	if err := swag.ReadJSON(raw, &dataFindScriptsOKBodyAO1); err != nil {
		return err
	}

	o.Data = dataFindScriptsOKBodyAO1.Data

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o FindScriptsOKBody) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	findScriptsOKBodyAO0, err := swag.WriteJSON(o.SearchResult)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, findScriptsOKBodyAO0)

	var dataFindScriptsOKBodyAO1 struct {
		Data []*models.ScriptOutput `json:"data"`
	}

	dataFindScriptsOKBodyAO1.Data = o.Data

	jsonDataFindScriptsOKBodyAO1, errFindScriptsOKBodyAO1 := swag.WriteJSON(dataFindScriptsOKBodyAO1)
	if errFindScriptsOKBodyAO1 != nil {
		return nil, errFindScriptsOKBodyAO1
	}
	_parts = append(_parts, jsonDataFindScriptsOKBodyAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this find scripts o k body
func (o *FindScriptsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with models.SearchResult
	if err := o.SearchResult.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *FindScriptsOKBody) validateData(formats strfmt.Registry) error {

	if swag.IsZero(o.Data) { // not required
		return nil
	}

	for i := 0; i < len(o.Data); i++ {
		if swag.IsZero(o.Data[i]) { // not required
			continue
		}

		if o.Data[i] != nil {
			if err := o.Data[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("findScriptsOK" + "." + "data" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *FindScriptsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *FindScriptsOKBody) UnmarshalBinary(b []byte) error {
	var res FindScriptsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
