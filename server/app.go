package server

import (
	"os"
	"runtime"

	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/MontFerret/ferret-server/pkg/history"
	"github.com/MontFerret/ferret-server/pkg/persistence"
	"github.com/MontFerret/ferret-server/pkg/projects"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret-server/server/controllers"
	"github.com/MontFerret/ferret-server/server/db"
	"github.com/MontFerret/ferret-server/server/http"
	"github.com/MontFerret/ferret-server/server/http/api/restapi/operations"
	"github.com/MontFerret/ferret/pkg/compiler"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

type Application struct {
	settings    Settings
	logger      *zerolog.Logger
	server      *http.Server
	db          *db.Manager
	projects    *projects.Service
	scripts     *scripts.Service
	history     *history.Service
	execution   *execution.Service
	persistence *persistence.Service
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

	historySrv, err := history.NewService(dbManager)

	if err != nil {
		return nil, errors.Wrap(err, "create history service")
	}

	perSrv, err := persistence.NewService(dbManager)

	if err != nil {
		return nil, errors.Wrap(err, "create persistence service")
	}

	fqlCompiler := compiler.New()

	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	queueSize := uint64(mem.Sys/MEGABYTE) * 2
	queue, err := execution.NewInMemoryQueue(uint64(queueSize))

	if err != nil {
		return nil, errors.Wrap(err, "create execution queue")
	}

	execSrv, err := execution.NewService(
		settings.Execution,
		&logger,
		dbManager,
		fqlCompiler,
		queue,
		history.ToStatusWriter(historySrv),
		history.ToLogWriter(historySrv),
		persistence.ToWriter(perSrv),
	)

	return &Application{
		settings,
		&logger,
		server,
		dbManager,
		projectsSrv,
		scriptsSrv,
		historySrv,
		execSrv,
		perSrv,
	}, nil
}

func (app *Application) Run() error {
	if err := app.configureProjectsController(); err != nil {
		return errors.Wrap(err, "configure project controller")
	}

	if err := app.configureScriptsController(); err != nil {
		return errors.Wrap(err, "configure scripts controller")
	}

	if err := app.configureExecutionController(); err != nil {
		return errors.Wrap(err, "configure execution controller")
	}

	if err := app.configurePersistenceController(); err != nil {
		return errors.Wrap(err, "configure persistence controller")
	}

	return app.server.Run()
}

func (app *Application) configureProjectsController() error {
	ctl, err := controllers.NewProjects(app.projects)

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
	ctl, err := controllers.NewScripts(app.scripts)

	if err != nil {
		return errors.Wrap(err, "new scripts controller")
	}

	app.server.API().CreateScriptHandler = operations.CreateScriptHandlerFunc(ctl.CreateScript)
	app.server.API().UpdateScriptHandler = operations.UpdateScriptHandlerFunc(ctl.UpdateScript)
	app.server.API().DeleteScriptHandler = operations.DeleteScriptHandlerFunc(ctl.DeleteScript)
	app.server.API().GetScriptHandler = operations.GetScriptHandlerFunc(ctl.GetScript)
	app.server.API().FindScriptsHandler = operations.FindScriptsHandlerFunc(ctl.FindScripts)

	return nil
}

func (app *Application) configureExecutionController() error {
	ctl, err := controllers.NewExecution(app.execution, app.history)

	if err != nil {
		return errors.Wrap(err, "new execution controller")
	}

	app.server.API().CreateExecutionHandler = operations.CreateExecutionHandlerFunc(ctl.Create)
	app.server.API().DeleteExecutionHandler = operations.DeleteExecutionHandlerFunc(ctl.Delete)
	app.server.API().FindExecutionsHandler = operations.FindExecutionsHandlerFunc(ctl.Find)
	app.server.API().GetExecutionHandler = operations.GetExecutionHandlerFunc(ctl.Get)

	return nil
}

func (app *Application) configurePersistenceController() error {
	ctl, err := controllers.NewPersistence(app.persistence)

	if err != nil {
		return errors.Wrap(err, "new persistence controller")
	}

	app.server.API().FindProjectDataHandler = operations.FindProjectDataHandlerFunc(ctl.FindAll)
	app.server.API().FindScriptDataHandler = operations.FindScriptDataHandlerFunc(ctl.Find)
	app.server.API().GetScriptDataHandler = operations.GetScriptDataHandlerFunc(ctl.Get)
	app.server.API().UpdateScriptDataHandler = operations.UpdateScriptDataHandlerFunc(ctl.Update)
	app.server.API().DeleteScriptDataHandler = operations.DeleteScriptDataHandlerFunc(ctl.Delete)

	return nil
}
