// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteProjectHandlerFunc turns a function with the right signature into a delete project handler
type DeleteProjectHandlerFunc func(DeleteProjectParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteProjectHandlerFunc) Handle(params DeleteProjectParams) middleware.Responder {
	return fn(params)
}

// DeleteProjectHandler interface for that can handle valid delete project params
type DeleteProjectHandler interface {
	Handle(DeleteProjectParams) middleware.Responder
}

// NewDeleteProject creates a new http.Handler for the delete project operation
func NewDeleteProject(ctx *middleware.Context, handler DeleteProjectHandler) *DeleteProject {
	return &DeleteProject{Context: ctx, Handler: handler}
}

/*DeleteProject swagger:route DELETE /projects/{projectID} deleteProject

Delete Project

*/
type DeleteProject struct {
	Context *middleware.Context
	Handler DeleteProjectHandler
}

func (o *DeleteProject) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteProjectParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
