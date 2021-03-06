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

// FindProjectDataHandlerFunc turns a function with the right signature into a find project data handler
type FindProjectDataHandlerFunc func(FindProjectDataParams) middleware.Responder

// Handle executing the request and returning a response
func (fn FindProjectDataHandlerFunc) Handle(params FindProjectDataParams) middleware.Responder {
	return fn(params)
}

// FindProjectDataHandler interface for that can handle valid find project data params
type FindProjectDataHandler interface {
	Handle(FindProjectDataParams) middleware.Responder
}

// NewFindProjectData creates a new http.Handler for the find project data operation
func NewFindProjectData(ctx *middleware.Context, handler FindProjectDataHandler) *FindProjectData {
	return &FindProjectData{Context: ctx, Handler: handler}
}

/*FindProjectData swagger:route GET /projects/{projectId}/data findProjectData

List Data

*/
type FindProjectData struct {
	Context *middleware.Context
	Handler FindProjectDataHandler
}

func (o *FindProjectData) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFindProjectDataParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// FindProjectDataOKBody find project data o k body
// swagger:model FindProjectDataOKBody
type FindProjectDataOKBody struct {
	models.SearchResult

	// data
	Data []*models.DataOutput `json:"data"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *FindProjectDataOKBody) UnmarshalJSON(raw []byte) error {
	// FindProjectDataOKBodyAO0
	var findProjectDataOKBodyAO0 models.SearchResult
	if err := swag.ReadJSON(raw, &findProjectDataOKBodyAO0); err != nil {
		return err
	}
	o.SearchResult = findProjectDataOKBodyAO0

	// FindProjectDataOKBodyAO1
	var dataFindProjectDataOKBodyAO1 struct {
		Data []*models.DataOutput `json:"data"`
	}
	if err := swag.ReadJSON(raw, &dataFindProjectDataOKBodyAO1); err != nil {
		return err
	}

	o.Data = dataFindProjectDataOKBodyAO1.Data

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o FindProjectDataOKBody) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	findProjectDataOKBodyAO0, err := swag.WriteJSON(o.SearchResult)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, findProjectDataOKBodyAO0)

	var dataFindProjectDataOKBodyAO1 struct {
		Data []*models.DataOutput `json:"data"`
	}

	dataFindProjectDataOKBodyAO1.Data = o.Data

	jsonDataFindProjectDataOKBodyAO1, errFindProjectDataOKBodyAO1 := swag.WriteJSON(dataFindProjectDataOKBodyAO1)
	if errFindProjectDataOKBodyAO1 != nil {
		return nil, errFindProjectDataOKBodyAO1
	}
	_parts = append(_parts, jsonDataFindProjectDataOKBodyAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this find project data o k body
func (o *FindProjectDataOKBody) Validate(formats strfmt.Registry) error {
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

func (o *FindProjectDataOKBody) validateData(formats strfmt.Registry) error {

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
					return ve.ValidateName("findProjectDataOK" + "." + "data" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *FindProjectDataOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *FindProjectDataOKBody) UnmarshalBinary(b []byte) error {
	var res FindProjectDataOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
