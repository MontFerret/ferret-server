package repositories

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/projects"
	"github.com/MontFerret/ferret-server/server/db/repositories/queries"
	"github.com/arangodb/go-driver"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"time"
)

type (
	ProjectRecord struct {
		Key string `json:"_key"`
		dal.Metadata
		projects.Project
	}

	ProjectRepository struct {
		client     driver.Client
		collection driver.Collection
	}
)

func NewProjectRepository(client driver.Client, db driver.Database, collectionName string) (*ProjectRepository, error) {
	ctx := context.Background()

	exists, err := db.CollectionExists(ctx, collectionName)

	if err != nil {
		return nil, err
	}

	var collection driver.Collection

	if exists {
		collection, err = db.Collection(ctx, collectionName)

		if err != nil {
			return nil, err
		}
	} else {
		collection, err = db.CreateCollection(ctx, collectionName, nil)

		if err != nil {
			return nil, err
		}
	}

	return &ProjectRepository{client, collection}, nil
}

func (repo *ProjectRepository) Get(ctx context.Context, id string) (*projects.ProjectEntity, error) {
	record := &ProjectRecord{}
	meta, err := repo.collection.ReadDocument(ctx, id, record)

	if err != nil {
		if driver.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return repo.fromRecord(meta, record), nil
}

func (repo *ProjectRepository) Find(ctx context.Context, q dal.Query) ([]*projects.ProjectEntity, error) {
	cursor, err := repo.collection.Database().Query(
		ctx,
		fmt.Sprintf(queries.FindAll, repo.collection.Name()),
		map[string]interface{}{
			"offset": q.Pagination.Size * (q.Pagination.Page - 1),
			"count":  q.Pagination.Size,
		},
	)

	if err != nil {
		return nil, err
	}

	result := make([]*projects.ProjectEntity, 0, q.Pagination.Size)

	defer cursor.Close()

	for cursor.HasMore() {
		record := &ProjectRecord{}

		meta, err := cursor.ReadDocument(ctx, record)

		if err != nil {
			return nil, err
		}

		result = append(result, repo.fromRecord(meta, record))
	}

	return result, nil
}

func (repo *ProjectRepository) Create(ctx context.Context, project *projects.Project) (*dal.Entity, error) {
	if project == nil {
		return nil, common.Error(common.ErrMissedArgument, "project")
	}

	if project.Name == "" {
		return nil, common.Error(common.ErrInvalidArgument, "empty name")
	}

	cursor, err := repo.collection.Database().Query(
		ctx,
		fmt.Sprintf(queries.FindOneByName, repo.collection.Name()),
		map[string]interface{}{
			"name": project.Name,
		},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	exists := cursor.HasMore()

	if exists {
		return nil, common.Errorf(common.ErrNotUnique, "project name %s", project.Name)
	}

	dbName, err := repo.generateID()

	if err != nil {
		return nil, errors.Wrapf(err, "generation new project id")
	}

	db, err := repo.client.CreateDatabase(ctx, dbName, nil)

	if err != nil {
		return nil, errors.Wrapf(err, "create new database with id %s", dbName)
	}

	createdAt := time.Now()

	meta, err := repo.collection.CreateDocument(ctx, &ProjectRecord{
		Key:     dbName,
		Project: *project,
		Metadata: dal.Metadata{
			CreatedAt: createdAt,
		},
	})

	if err != nil {
		db.Remove(ctx)

		return nil, err
	}

	return createdEntity(meta, createdAt), nil
}

func (repo *ProjectRepository) Update(ctx context.Context, project *projects.UpdateProject) (*dal.Entity, error) {
	if project == nil {
		return nil, common.Error(common.ErrMissedArgument, "project")
	}

	if project.ID == "" {
		return nil, common.Error(common.ErrInvalidOperation, "project model does not have ID")
	}

	updatedAt := time.Now()

	old := &projects.ProjectEntity{}

	updateCtx := driver.WithMergeObjects(driver.WithReturnOld(ctx, &old), false)

	meta, err := repo.collection.UpdateDocument(updateCtx, project.ID, &ProjectRecord{
		Project: project.Project,
		Metadata: dal.Metadata{
			UpdateAt: updatedAt,
		},
	})

	if err != nil {
		return nil, err
	}

	return updatedEntity(meta, old.CreatedAt, updatedAt), nil
}

func (repo *ProjectRepository) Delete(ctx context.Context, id string) error {
	db, err := repo.client.Database(ctx, id)

	if err != nil {
		return errors.Wrap(err, "find database")
	}

	if err := db.Remove(ctx); err != nil {
		return errors.Wrap(err, "remove database")
	}

	_, err = repo.collection.RemoveDocument(ctx, id)

	if err != nil {
		return errors.Wrap(err, "remove database record")
	}

	return nil
}

func (repo *ProjectRepository) generateID() (string, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	b := []rune(id.String())
	res := make([]rune, 0, len(b))
	res = append(res, []rune("f")...)

	return string(append(res, b[1:]...)), nil
}

func (repo *ProjectRepository) fromRecord(meta driver.DocumentMeta, record *ProjectRecord) *projects.ProjectEntity {
	return &projects.ProjectEntity{
		Entity: dal.Entity{
			ID:       meta.Key,
			Rev:      meta.Rev,
			Metadata: record.Metadata,
		},
		Project: projects.Project{
			Name:        record.Name,
			Description: record.Description,
		},
	}
}
