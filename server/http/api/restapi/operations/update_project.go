// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateProjectHandlerFunc turns a function with the right signature into a update project handler
type UpdateProjectHandlerFunc func(UpdateProjectParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateProjectHandlerFunc) Handle(params UpdateProjectParams) middleware.Responder {
	return fn(params)
}

// UpdateProjectHandler interface for that can handle valid update project params
type UpdateProjectHandler interface {
	Handle(UpdateProjectParams) middleware.Responder
}

// NewUpdateProject creates a new http.Handler for the update project operation
func NewUpdateProject(ctx *middleware.Context, handler UpdateProjectHandler) *UpdateProject {
	return &UpdateProject{Context: ctx, Handler: handler}
}

/*UpdateProject swagger:route PUT /projects/{projectID} updateProject

Update Project

*/
type UpdateProject struct {
	Context *middleware.Context
	Handler UpdateProjectHandler
}

func (o *UpdateProject) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateProjectParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
