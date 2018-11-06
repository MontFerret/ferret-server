package repositories

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret-server/server/db/repositories/queries"
	"github.com/arangodb/go-driver"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"time"
)

type (
	ScriptRecord struct {
		Key string `json:"_key"`
		dal.Metadata
		scripts.Script
	}

	ScriptRepository struct {
		collection driver.Collection
	}
)

func NewScriptRepository(db driver.Database, collectionName string) (*ScriptRepository, error) {
	if db == nil {
		return nil, common.Error(common.ErrMissedArgument, "database")
	}

	if collectionName == "" {
		return nil, common.Error(common.ErrMissedArgument, "collectionName")
	}

	ctx := context.Background()

	exists, err := db.CollectionExists(ctx, collectionName)

	if err != nil {
		return nil, errors.Wrap(err, "collection check")
	}

	var collection driver.Collection

	if exists {
		c, err := db.Collection(ctx, collectionName)

		if err != nil {
			return nil, errors.Wrap(err, "connect to collection")
		}

		collection = c
	} else {
		c, err := db.CreateCollection(ctx, collectionName, nil)

		if err != nil {
			return nil, errors.Wrap(err, "create new collection")
		}

		collection = c
	}

	return &ScriptRepository{collection}, nil
}

func (repo *ScriptRepository) Get(ctx context.Context, id string) (*scripts.ScriptEntity, error) {
	record := &ScriptRecord{}

	meta, err := repo.collection.ReadDocument(ctx, id, record)

	if err != nil {
		if driver.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return repo.fromRecord(meta, record), nil
}

func (repo *ScriptRepository) Find(ctx context.Context, query dal.Query) ([]*scripts.ScriptEntity, error) {
	cursor, err := repo.collection.Database().Query(
		ctx,
		fmt.Sprintf(queries.FindAll, repo.collection.Name()),
		map[string]interface{}{
			"offset": query.Pagination.Size * (query.Pagination.Page - 1),
			"count":  query.Pagination.Size,
		},
	)

	if err != nil {
		return nil, err
	}

	result := make([]*scripts.ScriptEntity, 0, query.Pagination.Size)

	defer cursor.Close()

	for cursor.HasMore() {
		record := &ScriptRecord{}

		meta, err := cursor.ReadDocument(ctx, record)

		if err != nil {
			return nil, err
		}

		result = append(result, repo.fromRecord(meta, record))
	}

	return result, nil
}

func (repo *ScriptRepository) Create(ctx context.Context, script *scripts.Script) (*dal.Entity, error) {
	if script == nil {
		return nil, common.Error(common.ErrMissedArgument, "script")
	}

	id, err := uuid.NewV4()

	if err != nil {
		return nil, errors.Wrap(err, "new id")
	}

	key := id.String()
	createdAt := time.Now()

	record := &ScriptRecord{
		Key: key,
		Metadata: dal.Metadata{
			CreatedAt: createdAt,
		},
		Script: *script,
	}

	meta, err := repo.collection.CreateDocument(ctx, record)

	if err != nil {
		return nil, errors.Wrap(err, "create script")
	}

	return createdEntity(meta, createdAt), nil
}

func (repo *ScriptRepository) Update(ctx context.Context, script *scripts.UpdateScript) (*dal.Entity, error) {
	if script == nil {
		return nil, common.Error(common.ErrMissedArgument, "script")
	}

	if script.ID == "" {
		return nil, common.Error(common.ErrInvalidOperation, "script model does not have ID")
	}

	updatedAt := time.Now()

	old := &scripts.ScriptEntity{}

	updateCtx := driver.WithMergeObjects(driver.WithReturnOld(ctx, old), false)

	meta, err := repo.collection.UpdateDocument(updateCtx, script.ID, &ScriptRecord{
		Script: script.Script,
		Metadata: dal.Metadata{
			UpdateAt: updatedAt,
		},
	})

	if err != nil {
		return nil, err
	}

	return updatedEntity(meta, old.CreatedAt, updatedAt), nil
}

func (repo *ScriptRepository) Delete(ctx context.Context, id string) error {
	_, err := repo.collection.RemoveDocument(ctx, id)

	if err != nil {
		return errors.Wrap(err, "remove script")
	}

	return nil
}

func (repo *ScriptRepository) fromRecord(meta driver.DocumentMeta, record *ScriptRecord) *scripts.ScriptEntity {
	return &scripts.ScriptEntity{
		Entity: dal.Entity{
			ID:       meta.Key,
			Rev:      meta.Rev,
			Metadata: record.Metadata,
		},
		Script: record.Script,
	}
}
