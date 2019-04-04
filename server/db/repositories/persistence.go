package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/persistence"
	"github.com/MontFerret/ferret-server/server/db/repositories/queries"
)

type (
	persistenceRecord struct {
		Key string `json:"_key"`
		dal.Metadata
		persistence.Record
	}

	PersistenceRepository struct {
		collection driver.Collection
	}
)

func NewPersistenceRepository(db driver.Database, collectionName string) (*PersistenceRepository, error) {
	ctx := context.Background()

	collection, err := initCollection(ctx, db, collectionName)

	if err != nil {
		return nil, err
	}

	err = ensureHashIndexes(ctx, collection, []hashIndex{
		{
			fields: []string{"script_id"},
			opts: &driver.EnsureHashIndexOptions{
				Sparse: true,
			},
		},
		{
			fields: []string{"job_id"},
			opts: &driver.EnsureHashIndexOptions{
				Sparse: true,
			},
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "create indexes")
	}

	return &PersistenceRepository{collection: collection}, nil
}

func (repo *PersistenceRepository) Create(ctx context.Context, record persistence.Record) (dal.Entity, error) {
	data := persistenceRecord{}
	data.Key = record.JobID
	data.Metadata.CreatedAt = time.Now()
	data.Record = record

	meta, err := repo.collection.CreateDocument(ctx, data)

	if err != nil {
		return dal.Entity{}, err
	}

	return dal.Entity{
		ID:  meta.Key,
		Rev: meta.Rev,
	}, nil
}

func (repo *PersistenceRepository) Update(ctx context.Context, record persistence.Record) (dal.Entity, error) {
	if record.JobID == "" {
		return dal.Entity{}, common.Error(common.ErrInvalidOperation, "data record does not have ID")
	}

	updatedAt := time.Now()

	old := persistence.RecordEntity{}
	updateCtx := driver.WithMergeObjects(driver.WithReturnOld(ctx, &old), false)
	meta, err := repo.collection.UpdateDocument(updateCtx, record.JobID, &persistenceRecord{
		Record: record,
		Metadata: dal.Metadata{
			UpdateAt: updatedAt,
		},
	})

	if err != nil {
		return dal.Entity{}, err
	}

	return updatedEntity(meta, old.CreatedAt, updatedAt), nil
}

func (repo *PersistenceRepository) Get(ctx context.Context, id string) (persistence.RecordEntity, error) {
	record := persistenceRecord{}
	meta, err := repo.collection.ReadDocument(ctx, id, &record)

	if err != nil {
		return persistence.RecordEntity{}, err
	}

	return repo.fromRecord(meta, record), nil
}

func (repo *PersistenceRepository) Delete(ctx context.Context, id string) error {
	_, err := repo.collection.RemoveDocument(ctx, id)

	if err != nil {
		return errors.Wrap(err, "remove data")
	}

	return nil
}

func (repo *PersistenceRepository) Find(ctx context.Context, q dal.Query) (persistence.QueryResult, error) {
	params := map[string]interface{}{}
	bindPaginationParams(params, q.Pagination)

	cursor, err := repo.collection.Database().Query(
		ctx,
		fmt.Sprintf(queries.FindAll, repo.collection.Name()),
		params,
	)

	if err != nil {
		return persistence.QueryResult{}, err
	}

	data := make([]persistence.RecordEntity, 0, q.Pagination.Count)

	defer cursor.Close()

	for cursor.HasMore() {
		record := persistenceRecord{}

		meta, err := cursor.ReadDocument(ctx, &record)

		if err != nil {
			return persistence.QueryResult{}, err
		}

		data = append(data, repo.fromRecord(meta, record))
	}

	result := persistence.QueryResult{
		QueryResult: dal.QueryResult{
			Count: uint64(len(data)),
		},
		Data: data,
	}

	length := len(data)

	if length > 0 {
		first := data[0]
		result.BeforeCursor = dal.NewCursor(first.CreatedAt)

		if length == int(q.Pagination.Count) {
			last := data[length-1]
			result.AfterCursor = dal.NewCursor(last.CreatedAt)
		}
	}

	return result, nil
}

func (repo *PersistenceRepository) FindByScriptID(ctx context.Context, scriptID string, q dal.Query) ([]persistence.RecordEntity, error) {
	cursor, err := repo.collection.Database().Query(
		ctx,
		fmt.Sprintf(queries.FindAllByScriptID, repo.collection.Name()),
		map[string]interface{}{
			"offset":    q.Pagination.Size * (q.Pagination.Page - 1),
			"count":     q.Pagination.Size,
			"script_id": scriptID,
		},
	)

	if err != nil {
		return nil, err
	}

	result := make([]persistence.RecordEntity, 0, q.Pagination.Size)

	defer cursor.Close()

	for cursor.HasMore() {
		record := persistenceRecord{}

		meta, err := cursor.ReadDocument(ctx, &record)

		if err != nil {
			return nil, err
		}

		result = append(result, repo.fromRecord(meta, record))
	}

	return result, nil
}

func (repo *PersistenceRepository) fromRecord(meta driver.DocumentMeta, record persistenceRecord) persistence.RecordEntity {
	return persistence.RecordEntity{
		Entity: dal.Entity{
			ID:       meta.Key,
			Rev:      meta.Rev,
			Metadata: record.Metadata,
		},
		Record: record.Record,
	}
}
