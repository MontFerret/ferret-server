package repositories

import (
	"context"
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
		return nil, errors.Wrap(err, "create hash indexes")
	}

	err = ensureSkipListIndexes(ctx, collection, []skipListIndex{
		{
			fields: []string{"created_at"},
			opts: &driver.EnsureSkipListIndexOptions{
				Unique: true,
			},
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "create skip list indexes")
	}

	return &PersistenceRepository{collection: collection}, nil
}

func (repo *PersistenceRepository) Create(ctx context.Context, record persistence.Record) (dal.Entity, error) {
	ts := time.Now()
	data := persistenceRecord{}
	data.Key = record.JobID
	data.Metadata.CreatedAt = &ts
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

	ts := time.Now()
	out := persistence.RecordEntity{}
	updateCtx := driver.WithMergeObjects(driver.WithReturnOld(ctx, &out), false)
	meta, err := repo.collection.UpdateDocument(updateCtx, record.JobID, &persistenceRecord{
		Record: record,
		Metadata: dal.Metadata{
			UpdateAt: &ts,
		},
	})

	if err != nil {
		return dal.Entity{}, err
	}

	return updatedEntity(meta, out.CreatedAt, &ts), nil
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
		queries.FindAll(repo.collection.Name()),
		params,
	)

	if err != nil {
		return persistence.QueryResult{}, err
	}

	data := make([]persistence.RecordEntity, 0, q.Pagination.Count+1)

	defer cursor.Close()

	for cursor.HasMore() {
		record := persistenceRecord{}

		meta, err := cursor.ReadDocument(ctx, &record)

		if err != nil {
			return persistence.QueryResult{}, err
		}

		data = append(data, repo.fromRecord(meta, record))
	}

	result := persistence.QueryResult{}
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

func (repo *PersistenceRepository) FindByScriptID(ctx context.Context, scriptID string, q dal.Query) (persistence.QueryResult, error) {
	params := map[string]interface{}{}
	bindPaginationParams(params, q.Pagination)
	params[queries.ParamFilterByScriptID] = scriptID

	cursor, err := repo.collection.Database().Query(
		ctx,
		queries.FindAllByScriptID(repo.collection.Name()),
		params,
	)

	if err != nil {
		return persistence.QueryResult{}, err
	}

	data := make([]persistence.RecordEntity, 0, q.Pagination.Count+1)

	defer cursor.Close()

	for cursor.HasMore() {
		record := persistenceRecord{}

		meta, err := cursor.ReadDocument(ctx, &record)

		if err != nil {
			return persistence.QueryResult{}, err
		}

		data = append(data, repo.fromRecord(meta, record))
	}

	result := persistence.QueryResult{}
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
