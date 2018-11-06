package db

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/projects"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret-server/server/db/repositories"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/pkg/errors"
	"sync"
)

const systemDBName = "ferret_server"
const projectsCollection = "projects"
const scriptsCollection = "scripts"

type Manager struct {
	client           driver.Client
	systemDB         driver.Database
	projects         projects.Repository
	projectDatabases sync.Map
	scripts          sync.Map
}

func New(settings Settings) (*Manager, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: settings.Endpoints,
	})

	if err != nil {
		return nil, errors.Wrap(err, "open connection")
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})

	if err != nil {
		return nil, errors.Wrap(err, "create client")
	}

	ctx := context.Background()

	exists, err := client.DatabaseExists(ctx, systemDBName)

	if err != nil {
		return nil, errors.Wrap(err, "system database check")
	}

	var db driver.Database

	if !exists {
		db, err = client.CreateDatabase(ctx, systemDBName, nil)

		if err != nil {
			return nil, errors.Wrapf(err, "create system database %s", systemDBName)
		}
	} else {
		db, err = client.Database(ctx, systemDBName)

		if err != nil {
			return nil, errors.Wrapf(err, "open system database %s", systemDBName)
		}
	}

	proj, err := repositories.NewProjectRepository(client, db, projectsCollection)

	if err != nil {
		return nil, err
	}

	manager := new(Manager)
	manager.client = client
	manager.systemDB = db
	manager.projects = proj
	manager.scripts = sync.Map{}

	return manager, nil
}

func (manager *Manager) GetProjectsRepository() (projects.Repository, error) {
	return manager.projects, nil
}

func (manager *Manager) GetScriptsRepository(projectID string) (scripts.Repository, error) {
	repo, exists := manager.scripts.Load(projectID)

	if !exists {
		db, exists := manager.projectDatabases.Load(projectID)

		if !exists {
			d, err := manager.client.Database(context.Background(), projectID)

			if err != nil {
				return nil, errors.Wrap(err, "connect to database")
			}

			manager.projectDatabases.Store(projectID, d)

			db = d
		}

		r, err := repositories.NewScriptRepository(db.(driver.Database), scriptsCollection)

		if err != nil {
			return nil, err
		}

		manager.scripts.Store(projectID, r)

		repo = r
	}

	return repo.(scripts.Repository), nil
}
