package repositories

import (
	"context"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret-server/server/db/repositories/queries"
)

type (
	scriptRecord struct {
		Key string `json:"_key"`
		dal.Metadata
		scripts.Script
	}

	ScriptRepository struct {
		collection driver.Collection
	}
)

func NewScriptRepository(db driver.Database, collectionName string) (*ScriptRepository, error) {
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

	return &ScriptRepository{collection}, nil
}

func (repo *ScriptRepository) Get(ctx context.Context, id string) (scripts.ScriptEntity, error) {
	record := scriptRecord{}

	meta, err := repo.collection.ReadDocument(ctx, id, &record)

	if err != nil {
		if driver.IsNotFound(err) {
			return scripts.ScriptEntity{}, common.ErrNotFound
		}

		return scripts.ScriptEntity{}, err
	}

	return repo.fromRecord(meta, record), nil
}

func (repo *ScriptRepository) Find(ctx context.Context, q dal.Query) (scripts.QueryResult, error) {
	params := map[string]interface{}{}
	bindPaginationParams(params, q.Pagination)

	cursor, err := repo.collection.Database().Query(
		ctx,
		queries.FindAll(repo.collection.Name()),
		params,
	)

	if err != nil {
		return scripts.QueryResult{}, err
	}

	data := make([]scripts.ScriptEntity, 0, q.Pagination.Count+1)

	defer cursor.Close()

	for cursor.HasMore() {
		record := scriptRecord{}

		meta, err := cursor.ReadDocument(ctx, &record)

		if err != nil {
			return scripts.QueryResult{}, err
		}

		data = append(data, repo.fromRecord(meta, record))
	}

	result := scripts.QueryResult{}
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

func (repo *ScriptRepository) Create(ctx context.Context, script scripts.Script) (dal.Entity, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return dal.Entity{}, errors.Wrap(err, "new id")
	}

	key := id.String()
	ts := time.Now()

	record := scriptRecord{
		Key: key,
		Metadata: dal.Metadata{
			CreatedAt: &ts,
		},
		Script: script,
	}

	meta, err := repo.collection.CreateDocument(ctx, record)

	if err != nil {
		return dal.Entity{}, errors.Wrap(err, "create script")
	}

	return createdEntity(meta, &ts), nil
}

func (repo *ScriptRepository) Update(ctx context.Context, script scripts.UpdateScript) (dal.Entity, error) {
	if script.ID == "" {
		return dal.Entity{}, common.Error(common.ErrInvalidOperation, "script model does not have ID")
	}

	ts := time.Now()
	out := &scripts.ScriptEntity{}
	updateCtx := driver.WithMergeObjects(driver.WithReturnOld(ctx, out), false)

	meta, err := repo.collection.UpdateDocument(updateCtx, script.ID, &scriptRecord{
		Script: script.Script,
		Metadata: dal.Metadata{
			UpdateAt: &ts,
		},
	})

	if err != nil {
		return dal.Entity{}, err
	}

	return updatedEntity(meta, out.CreatedAt, &ts), nil
}

func (repo *ScriptRepository) Delete(ctx context.Context, id string) error {
	_, err := repo.collection.RemoveDocument(ctx, id)

	if err != nil {
		return errors.Wrap(err, "remove script")
	}

	return nil
}

func (repo *ScriptRepository) fromRecord(meta driver.DocumentMeta, record scriptRecord) scripts.ScriptEntity {
	return scripts.ScriptEntity{
		Entity: dal.Entity{
			ID:       meta.Key,
			Rev:      meta.Rev,
			Metadata: record.Metadata,
		},
		Script: record.Script,
	}
}
