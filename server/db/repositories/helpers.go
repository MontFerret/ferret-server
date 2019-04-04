package repositories

import (
	"bytes"
	"context"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/server/db/repositories/queries"
)

type (
	hashIndex struct {
		fields []string
		opts   *driver.EnsureHashIndexOptions
	}

	skipListIndex struct {
		fields []string
		opts   *driver.EnsureSkipListIndexOptions
	}
)

func createdEntity(meta driver.DocumentMeta, time time.Time) dal.Entity {
	return dal.Entity{
		ID:  meta.Key,
		Rev: meta.Rev,
		Metadata: dal.Metadata{
			CreatedAt: time,
		},
	}
}

func updatedEntity(meta driver.DocumentMeta, createdAt, updatedAt time.Time) dal.Entity {
	res := createdEntity(meta, createdAt)
	res.UpdateAt = updatedAt

	return res
}

func initCollection(ctx context.Context, db driver.Database, collectionName string) (driver.Collection, error) {
	if db == nil {
		return nil, common.Error(common.ErrMissedArgument, "database")
	}

	if collectionName == "" {
		return nil, common.Error(common.ErrMissedArgument, "collectionName")
	}

	exists, err := db.CollectionExists(ctx, collectionName)

	if err != nil {
		return nil, errors.Wrap(err, "collection check")
	}

	if exists {
		c, err := db.Collection(ctx, collectionName)

		if err != nil {
			return nil, errors.Wrap(err, "connect to collection")
		}

		return c, nil
	}

	c, err := db.CreateCollection(ctx, collectionName, nil)

	if err != nil {
		return nil, errors.Wrap(err, "create new collection")
	}

	return c, nil
}

func ensureHashIndexes(ctx context.Context, collection driver.Collection, indexes []hashIndex) error {
	for _, i := range indexes {
		_, _, err := collection.EnsureHashIndex(ctx, i.fields, i.opts)

		if err != nil {
			return err
		}
	}

	return nil
}

func ensureSkipListIndexes(ctx context.Context, collection driver.Collection, indexes []skipListIndex) error {
	for _, i := range indexes {
		_, _, err := collection.EnsureSkipListIndex(ctx, i.fields, i.opts)

		if err != nil {
			return err
		}
	}

	return nil
}

func compileQuery(collectionName string, q dal.Query) dal.CompiledQuery {
	var qs bytes.Buffer
	params := map[string]interface{}{}

	varName := "i"
	qs.WriteString("FOR ")
	qs.WriteString(varName)
	qs.WriteString(" IN ")
	qs.WriteString(collectionName)

	if q.Filtering.Fields != nil && len(q.Filtering.Fields) > 0 {
		qs.WriteString("\n")
		qs.WriteString("FILTER ")

		lastIndex := len(q.Filtering.Fields) - 1
		for i, f := range q.Filtering.Fields {
			paramName := f.Name
			qs.WriteString(varName)
			qs.WriteString(".")
			qs.WriteString(f.Name)
			qs.WriteString(" ")
			qs.WriteString(f.Comparator.String())
			qs.WriteString(" ")
			qs.WriteString("@")
			qs.WriteString(paramName)

			params[paramName] = f.Value

			if i != lastIndex {
				qs.WriteString(" ")
				qs.WriteString(q.Filtering.Operator.String())
				qs.WriteString(" ")
			}
		}
	}

	if q.Pagination.Cursor.IsEmpty() == false {
		qs.WriteString("\n")
		qs.WriteString("FILTER ")
		qs.WriteString(varName)
		qs.WriteString(".created_at < @")
		qs.WriteString(queries.ParamPageCursor)

		params[queries.ParamPageCursor] = q.Pagination.Cursor
	}

	if q.Pagination.Count > 0 {
		qs.WriteString("\n")
		qs.WriteString("LIMIT @")
		qs.WriteString(queries.ParamPageCount)

		params[queries.ParamPageCount] = q.Pagination.Count
	}

	qs.WriteString("\n")
	qs.WriteString("RETURN i")

	return dal.CompiledQuery{
		String: qs.String(),
		Params: params,
	}
}

func bindPaginationParams(params map[string]interface{}, p dal.Pagination) {
	if p.Cursor.IsEmpty() == false {
		params[queries.ParamPageCursor] = p.Cursor
	} else {
		params[queries.ParamPageCursor] = nil
	}

	params[queries.ParamPageCount] = p.Count
}
