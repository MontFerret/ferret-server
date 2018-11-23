package repositories

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"
	"time"
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

func initCollection(db driver.Database, collectionName string) (driver.Collection, error) {
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
