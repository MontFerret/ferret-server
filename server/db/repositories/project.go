package repositories

import (
	"context"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/projects"
	"github.com/MontFerret/ferret-server/server/db/repositories/queries"
)

type (
	projectRecord struct {
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

	collection, err := initCollection(ctx, db, collectionName)

	if err != nil {
		return nil, err
	}

	err = ensureSkipListIndexes(ctx, collection, []skipListIndex{
		{
			fields: []string{"name"},
			opts: &driver.EnsureSkipListIndexOptions{
				Unique: true,
			},
		},
		{
			fields: []string{"created_at"},
			opts: &driver.EnsureSkipListIndexOptions{
				Unique: true,
			},
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "create skiplist indexes")
	}

	return &ProjectRepository{client, collection}, nil
}

func (repo *ProjectRepository) Get(ctx context.Context, id string) (projects.ProjectEntity, error) {
	record := projectRecord{}
	meta, err := repo.collection.ReadDocument(ctx, id, &record)

	if err != nil {
		if driver.IsNotFound(err) {
			return projects.ProjectEntity{}, common.ErrNotFound
		}

		return projects.ProjectEntity{}, err
	}

	return repo.fromRecord(meta, record), nil
}

func (repo *ProjectRepository) Find(ctx context.Context, q dal.Query) (projects.QueryResult, error) {
	params := map[string]interface{}{}
	bindPaginationParams(params, q.Pagination)

	cursor, err := repo.collection.Database().Query(
		ctx,
		queries.FindAll(repo.collection.Name()),
		params,
	)

	if err != nil {
		return projects.QueryResult{}, err
	}

	data := make([]projects.ProjectEntity, 0, q.Pagination.Count+1)

	defer cursor.Close()

	for cursor.HasMore() {
		record := projectRecord{}

		meta, err := cursor.ReadDocument(ctx, &record)

		if err != nil {
			return projects.QueryResult{}, err
		}

		data = append(data, repo.fromRecord(meta, record))
	}

	result := projects.QueryResult{}
	length := len(data)
	result.QueryResult = createPaginationResult(q.Pagination, length)

	if length > 0 {
		if length >= int(q.Pagination.Count) {
			result.Data = data[:q.Pagination.Count]
		} else {
			result.Data = data
		}
	}

	return result, nil
}

func (repo *ProjectRepository) Create(ctx context.Context, project projects.Project) (dal.Entity, error) {
	if project.Name == "" {
		return dal.Entity{}, common.Error(common.ErrInvalidArgument, "empty name")
	}

	cursor, err := repo.collection.Database().Query(
		ctx,
		queries.FindOneByName(repo.collection.Name()),
		map[string]interface{}{
			"name": project.Name,
		},
	)

	if err != nil {
		return dal.Entity{}, err
	}

	defer cursor.Close()

	exists := cursor.HasMore()

	if exists {
		return dal.Entity{}, common.Errorf(common.ErrNotUnique, "project name %s", project.Name)
	}

	dbName, err := repo.generateID()

	if err != nil {
		return dal.Entity{}, errors.Wrapf(err, "generation new project id")
	}

	db, err := repo.client.CreateDatabase(ctx, dbName, nil)

	if err != nil {
		return dal.Entity{}, errors.Wrapf(err, "create new database with id %s", dbName)
	}

	ts := time.Now()
	meta, err := repo.collection.CreateDocument(ctx, &projectRecord{
		Key:     dbName,
		Project: project,
		Metadata: dal.Metadata{
			CreatedAt: &ts,
		},
	})

	if err != nil {
		db.Remove(ctx)

		return dal.Entity{}, err
	}

	return createdEntity(meta, &ts), nil
}

func (repo *ProjectRepository) Update(ctx context.Context, project projects.UpdateProject) (dal.Entity, error) {
	if project.ID == "" {
		return dal.Entity{}, common.Error(common.ErrInvalidOperation, "project model does not have ID")
	}

	ts := time.Now()
	out := &projects.ProjectEntity{}
	updateCtx := driver.WithMergeObjects(driver.WithReturnOld(ctx, &out), false)

	meta, err := repo.collection.UpdateDocument(updateCtx, project.ID, &projectRecord{
		Project: project.Project,
		Metadata: dal.Metadata{
			UpdateAt: &ts,
		},
	})

	if err != nil {
		return dal.Entity{}, err
	}

	return updatedEntity(meta, out.CreatedAt, &ts), nil
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

func (repo *ProjectRepository) fromRecord(meta driver.DocumentMeta, record projectRecord) projects.ProjectEntity {
	return projects.ProjectEntity{
		Entity: dal.Entity{
			ID:       meta.Key,
			Rev:      meta.Rev,
			Metadata: record.Metadata,
		},
		Project: record.Project,
	}
}
