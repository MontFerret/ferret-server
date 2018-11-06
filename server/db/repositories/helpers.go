package repositories

import (
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/arangodb/go-driver"
	"time"
)

func createdEntity(meta driver.DocumentMeta, time time.Time) *dal.Entity {
	return &dal.Entity{
		ID:  meta.Key,
		Rev: meta.Rev,
		Metadata: dal.Metadata{
			CreatedAt: time,
		},
	}
}

func updatedEntity(meta driver.DocumentMeta, createdAt, updatedAt time.Time) *dal.Entity {
	res := createdEntity(meta, createdAt)
	res.UpdateAt = updatedAt

	return res
}
