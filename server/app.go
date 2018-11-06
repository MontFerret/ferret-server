package server

import (
	"github.com/MontFerret/ferret-server/pkg/projects"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret-server/server/controllers"
	"github.com/MontFerret/ferret-server/server/db"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
)

type Application struct {
	settings Settings
	logger   *zerolog.Logger
	server   *http.Server
	db       *db.Manager
	projects *projects.Service
	scripts  *scripts.Service
}

func New(settings Settings) (*Application, error) {
	logger := zerolog.New(os.Stdout).
		With().
		Str("name", settings.Name).
		Str("version", settings.Version).
		Logger()

	server, err := http.New(&logger, settings.HTTP)

	if err != nil {
		return nil, errors.Wrap(err, "create server")
	}

	dbManager, err := db.New(settings.Database)

	if err != nil {
		return nil, errors.Wrap(err, "connect to db")
	}

	projectsSrv, err := projects.NewService(dbManager)

	if err != nil {
		return nil, errors.Wrap(err, "create projects service")
	}

	scriptsSrv, err := scripts.NewService(dbManager)

	if err != nil {
		return nil, errors.Wrap(err, "create scripts service")
	}

	return &Application{
		settings,
		&logger,
		server,
		dbManager,
		projectsSrv,
		scriptsSrv,
	}, nil
}

func (app *Application) Run() error {
	if err := app.configureProjectsController(); err != nil {
		return errors.Wrap(err, "configure project controller")
	}

	if err := app.configureScriptsController(); err != nil {
		return errors.Wrap(err, "configure scripts controller")
	}

	return app.server.Run()
}

func (app *Application) configureProjectsController() error {
	ctl, err := controllers.NewProjectsController(app.projects)

	if err != nil {
		return errors.Wrap(err, "new projects controller")
	}

	app.server.API().CreateProjectHandler = operations.CreateProjectHandlerFunc(ctl.CreateProject)
	app.server.API().UpdateProjectHandler = operations.UpdateProjectHandlerFunc(ctl.UpdateProject)
	app.server.API().DeleteProjectHandler = operations.DeleteProjectHandlerFunc(ctl.DeleteProject)
	app.server.API().GetProjectHandler = operations.GetProjectHandlerFunc(ctl.GetProject)
	app.server.API().FindProjectsHandler = operations.FindProjectsHandlerFunc(ctl.FindProjects)

	return nil
}

func (app *Application) configureScriptsController() error {
	ctl, err := controllers.NewScriptsController(app.scripts)

	if err != nil {
		return errors.Wrap(err, "new scripts controller")
	}

	app.server.API().CreateScriptHandler = operations.CreateScriptHandlerFunc(ctl.CreateScript)
	app.server.API().UpdateScriptHandler = operations.UpdateScriptHandlerFunc(ctl.UpdateScript)
	app.server.API().DeleteScriptHandler = operations.DeleteScriptHandlerFunc(ctl.DeleteScript)
	app.server.API().GetScriptHandler = operations.GetScriptHandlerFunc(ctl.GetScripts)
	app.server.API().FindScriptsHandler = operations.FindScriptsHandlerFunc(ctl.FindScripts)

	return nil
}
