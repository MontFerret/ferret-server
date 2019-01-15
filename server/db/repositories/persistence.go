package repositories

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/persistence"
	"github.com/MontFerret/ferret-server/server/db/repositories/queries"
	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"
	"time"
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
	collection, err := initCollection(db, collectionName)

	if err != nil {
		return nil, err
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

func (repo *PersistenceRepository) Find(ctx context.Context, q dal.Query) ([]persistence.RecordEntity, error) {
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
