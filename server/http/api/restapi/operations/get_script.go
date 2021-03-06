// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetScriptHandlerFunc turns a function with the right signature into a get script handler
type GetScriptHandlerFunc func(GetScriptParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetScriptHandlerFunc) Handle(params GetScriptParams) middleware.Responder {
	return fn(params)
}

// GetScriptHandler interface for that can handle valid get script params
type GetScriptHandler interface {
	Handle(GetScriptParams) middleware.Responder
}

// NewGetScript creates a new http.Handler for the get script operation
func NewGetScript(ctx *middleware.Context, handler GetScriptHandler) *GetScript {
	return &GetScript{Context: ctx, Handler: handler}
}

/*GetScript swagger:route GET /projects/{projectId}/scripts/{scriptId} getScript

Get Script

*/
type GetScript struct {
	Context *middleware.Context
	Handler GetScriptHandler
}

func (o *GetScript) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetScriptParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
