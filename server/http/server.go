package http

import (
	"github.com/MontFerret/ferret-server/server/http/api/restapi"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/MontFerret/ferret-server/server/logging"
	"github.com/codegangsta/negroni"
	oaerrors "github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/pilu/xrequestid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Server struct {
	settings Settings
	logger   *zerolog.Logger
	api      *operations.FerretServerAPI
	engine   *restapi.Server
}

func New(logger *zerolog.Logger, settings Settings) (*Server, error) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)

	if err != nil {
		return nil, errors.Wrap(err, "load spec")
	}

	api := operations.NewFerretServerAPI(swaggerSpec)
	api.Logger = func(s string, i ...interface{}) {
		logger.Info().Timestamp().Msgf(s, i...)
	}
	engine := restapi.NewServer(api)
	engine.Port = int(settings.Port)
	api.ServeError = oaerrors.ServeError
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.ServerShutdown = func() {}

	return &Server{settings, logger, api, engine}, nil
}

func (server *Server) API() *operations.FerretServerAPI {
	return server.api
}

func (server *Server) Run() error {
	n := negroni.New()
	n.Use(xrequestid.New(4))
	n.Use(logging.NewMiddleware(server.logger, xrequestid.DefaultHeaderKey))
	n.Use(NewRecovery())
	n.UseHandler(server.api.Serve(nil))

	server.engine.SetHandler(n)

	err := server.engine.Serve()
	defer server.engine.Shutdown()

	if err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	return nil
}
