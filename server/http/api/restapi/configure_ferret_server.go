// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
)

//go:generate swagger generate server --target ../../api --name FerretServer --spec ../../../../api/api.oas2.json --exclude-main

func configureFlags(api *operations.FerretServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.FerretServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.CreateExecutionHandler == nil {
		api.CreateExecutionHandler = operations.CreateExecutionHandlerFunc(func(params operations.CreateExecutionParams) middleware.Responder {
			return middleware.NotImplemented("operation .CreateExecution has not yet been implemented")
		})
	}
	if api.CreateProjectHandler == nil {
		api.CreateProjectHandler = operations.CreateProjectHandlerFunc(func(params operations.CreateProjectParams) middleware.Responder {
			return middleware.NotImplemented("operation .CreateProject has not yet been implemented")
		})
	}
	if api.CreateScriptHandler == nil {
		api.CreateScriptHandler = operations.CreateScriptHandlerFunc(func(params operations.CreateScriptParams) middleware.Responder {
			return middleware.NotImplemented("operation .CreateScript has not yet been implemented")
		})
	}
	if api.DeleteExecutionHandler == nil {
		api.DeleteExecutionHandler = operations.DeleteExecutionHandlerFunc(func(params operations.DeleteExecutionParams) middleware.Responder {
			return middleware.NotImplemented("operation .DeleteExecution has not yet been implemented")
		})
	}
	if api.DeleteProjectHandler == nil {
		api.DeleteProjectHandler = operations.DeleteProjectHandlerFunc(func(params operations.DeleteProjectParams) middleware.Responder {
			return middleware.NotImplemented("operation .DeleteProject has not yet been implemented")
		})
	}
	if api.DeleteScriptHandler == nil {
		api.DeleteScriptHandler = operations.DeleteScriptHandlerFunc(func(params operations.DeleteScriptParams) middleware.Responder {
			return middleware.NotImplemented("operation .DeleteScript has not yet been implemented")
		})
	}
	if api.DeleteScriptDataHandler == nil {
		api.DeleteScriptDataHandler = operations.DeleteScriptDataHandlerFunc(func(params operations.DeleteScriptDataParams) middleware.Responder {
			return middleware.NotImplemented("operation .DeleteScriptData has not yet been implemented")
		})
	}
	if api.FindExecutionsHandler == nil {
		api.FindExecutionsHandler = operations.FindExecutionsHandlerFunc(func(params operations.FindExecutionsParams) middleware.Responder {
			return middleware.NotImplemented("operation .FindExecutions has not yet been implemented")
		})
	}
	if api.FindProjectDataHandler == nil {
		api.FindProjectDataHandler = operations.FindProjectDataHandlerFunc(func(params operations.FindProjectDataParams) middleware.Responder {
			return middleware.NotImplemented("operation .FindProjectData has not yet been implemented")
		})
	}
	if api.FindProjectsHandler == nil {
		api.FindProjectsHandler = operations.FindProjectsHandlerFunc(func(params operations.FindProjectsParams) middleware.Responder {
			return middleware.NotImplemented("operation .FindProjects has not yet been implemented")
		})
	}
	if api.FindScriptDataHandler == nil {
		api.FindScriptDataHandler = operations.FindScriptDataHandlerFunc(func(params operations.FindScriptDataParams) middleware.Responder {
			return middleware.NotImplemented("operation .FindScriptData has not yet been implemented")
		})
	}
	if api.FindScriptsHandler == nil {
		api.FindScriptsHandler = operations.FindScriptsHandlerFunc(func(params operations.FindScriptsParams) middleware.Responder {
			return middleware.NotImplemented("operation .FindScripts has not yet been implemented")
		})
	}
	if api.GetExecutionHandler == nil {
		api.GetExecutionHandler = operations.GetExecutionHandlerFunc(func(params operations.GetExecutionParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetExecution has not yet been implemented")
		})
	}
	if api.GetProjectHandler == nil {
		api.GetProjectHandler = operations.GetProjectHandlerFunc(func(params operations.GetProjectParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetProject has not yet been implemented")
		})
	}
	if api.GetScriptHandler == nil {
		api.GetScriptHandler = operations.GetScriptHandlerFunc(func(params operations.GetScriptParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetScript has not yet been implemented")
		})
	}
	if api.GetScriptDataHandler == nil {
		api.GetScriptDataHandler = operations.GetScriptDataHandlerFunc(func(params operations.GetScriptDataParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetScriptData has not yet been implemented")
		})
	}
	if api.UpdateProjectHandler == nil {
		api.UpdateProjectHandler = operations.UpdateProjectHandlerFunc(func(params operations.UpdateProjectParams) middleware.Responder {
			return middleware.NotImplemented("operation .UpdateProject has not yet been implemented")
		})
	}
	if api.UpdateScriptHandler == nil {
		api.UpdateScriptHandler = operations.UpdateScriptHandlerFunc(func(params operations.UpdateScriptParams) middleware.Responder {
			return middleware.NotImplemented("operation .UpdateScript has not yet been implemented")
		})
	}
	if api.UpdateScriptDataHandler == nil {
		api.UpdateScriptDataHandler = operations.UpdateScriptDataHandlerFunc(func(params operations.UpdateScriptDataParams) middleware.Responder {
			return middleware.NotImplemented("operation .UpdateScriptData has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
