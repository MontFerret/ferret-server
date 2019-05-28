package repositories

import (
	"context"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/history"
)

type (
	historyRecord struct {
		Key string `json:"_key"`
		dal.Metadata
		history.Record
	}

	HistoryRepository struct {
		collection driver.Collection
	}
)

func NewHistoryRepository(db driver.Database, collectionName string) (*HistoryRepository, error) {
	ctx := context.Background()
	collection, err := initCollection(ctx, db, collectionName)

	if err != nil {
		return nil, err
	}

	err = ensureHashIndexes(ctx, collection, []hashIndex{
		{
			fields: []string{"script_id"},
		},
		{
			fields: []string{"job_id"},
		},
		{
			fields: []string{"status"},
		},
		{
			fields: []string{"cause"},
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "create hash indexes")
	}

	err = ensureSkipListIndexes(ctx, collection, []skipListIndex{
		{
			fields: []string{"started_at"},
		},
		{
			fields: []string{"ended_at"},
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "create skiplist indexes")
	}

	return &HistoryRepository{collection}, nil
}

func (repo *HistoryRepository) Create(ctx context.Context, entry history.Record) (dal.Entity, error) {
	record := historyRecord{}
	ts := time.Now()
	record.Key = entry.JobID
	record.Metadata.CreatedAt = &ts
	record.Record = entry

	meta, err := repo.collection.CreateDocument(ctx, record)

	if err != nil {
		return dal.Entity{}, err
	}

	return dal.Entity{
		ID:  meta.Key,
		Rev: meta.Rev,
	}, nil
}

func (repo *HistoryRepository) Update(ctx context.Context, entry history.Record) (dal.Entity, error) {
	if entry.JobID == "" {
		return dal.Entity{}, common.Error(common.ErrInvalidOperation, "project model does not have ID")
	}

	ts := time.Now()
	out := history.RecordEntity{}
	updateCtx := driver.WithMergeObjects(driver.WithReturnOld(ctx, &out), false)
	meta, err := repo.collection.UpdateDocument(updateCtx, entry.JobID, &historyRecord{
		Record: entry,
		Metadata: dal.Metadata{
			UpdateAt: &ts,
		},
	})

	if err != nil {
		return dal.Entity{}, err
	}

	return updatedEntity(meta, out.CreatedAt, &ts), nil
}

func (repo *HistoryRepository) Get(ctx context.Context, jobID string) (history.RecordEntity, error) {
	record := historyRecord{}
	meta, err := repo.collection.ReadDocument(ctx, jobID, &record)

	if err != nil {
		return history.RecordEntity{}, err
	}

	return repo.fromRecord(meta, record), nil
}

func (repo *HistoryRepository) Find(ctx context.Context, q dal.Query) (history.QueryResult, error) {
	cq := compileQuery(repo.collection.Name(), q)

	cursor, err := repo.collection.Database().Query(
		ctx,
		cq.String,
		cq.Params,
	)

	if err != nil {
		return history.QueryResult{}, err
	}

	data := make([]history.RecordEntity, 0, q.Pagination.Count+1)

	defer cursor.Close()

	for cursor.HasMore() {
		record := historyRecord{}

		meta, err := cursor.ReadDocument(ctx, &record)

		if err != nil {
			return history.QueryResult{}, err
		}

		data = append(data, repo.fromRecord(meta, record))
	}

	result := history.QueryResult{}
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

func (repo *HistoryRepository) fromRecord(meta driver.DocumentMeta, record historyRecord) history.RecordEntity {
	return history.RecordEntity{
		Entity: dal.Entity{
			ID:       meta.Key,
			Rev:      meta.Rev,
			Metadata: record.Metadata,
		},
		Record: record.Record,
	}
}
