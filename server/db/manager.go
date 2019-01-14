package db

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/history"
	"github.com/MontFerret/ferret-server/pkg/persistence"
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
const historyCollection = "history"
const dataCollection = "data"

type (
	factory func(driver.Database, string) (interface{}, error)

	Manager struct {
		client    driver.Client
		systemDB  driver.Database
		databases sync.Map
		projects  projects.Repository
		scripts   sync.Map
		histories sync.Map
		data      sync.Map
	}
)

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
	manager.data = sync.Map{}

	return manager, nil
}

func (manager *Manager) GetProjectsRepository() (projects.Repository, error) {
	return manager.projects, nil
}

func (manager *Manager) GetScriptsRepository(projectID string) (scripts.Repository, error) {
	repo, err := manager.resolveRepo(projectID, scriptsCollection, manager.scripts, func(db driver.Database, name string) (interface{}, error) {
		return repositories.NewScriptRepository(db, name)
	})

	if err != nil {
		return nil, err
	}

	return repo.(scripts.Repository), nil
}

func (manager *Manager) GetHistoryRepository(projectID string) (history.Repository, error) {
	repo, err := manager.resolveRepo(projectID, historyCollection, manager.histories, func(db driver.Database, name string) (interface{}, error) {
		return repositories.NewHistoryRepository(db, name)
	})

	if err != nil {
		return nil, err
	}

	return repo.(history.Repository), nil
}

func (manager *Manager) GetDataRepository(projectID string) (persistence.Repository, error) {
	repo, err := manager.resolveRepo(projectID, dataCollection, manager.data, func(db driver.Database, name string) (interface{}, error) {
		return repositories.NewDataRepository(db, name)
	})

	if err != nil {
		return nil, err
	}

	return repo.(persistence.Repository), nil
}

func (manager *Manager) resolveRepo(projectID string, collectionName string, cache sync.Map, f factory) (interface{}, error) {
	repo, exists := cache.Load(projectID)

	if !exists {
		db, err := manager.resolveDB(projectID)

		if err != nil {
			return nil, err
		}

		r, err := f(db, collectionName)

		if err != nil {
			return nil, err
		}

		cache.Store(projectID, r)

		repo = r
	}

	return repo, nil
}

func (manager *Manager) resolveDB(projectID string) (driver.Database, error) {
	// get a db connection from the cache
	db, exists := manager.databases.Load(projectID)

	if !exists {
		// connect to the project DB
		d, err := manager.client.Database(context.Background(), projectID)

		if err != nil {
			return nil, errors.Wrap(err, "connect to database")
		}

		// cache it
		manager.databases.Store(projectID, d)

		db = d
	}

	return db.(driver.Database), nil
}
