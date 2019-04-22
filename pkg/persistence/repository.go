package persistence

import (
	"context"

	"github.com/MontFerret/ferret-server/pkg/common/dal"
)

type Repository interface {
	Create(ctx context.Context, record Record) (dal.Entity, error)

	Update(ctx context.Context, record Record) (dal.Entity, error)

	Get(ctx context.Context, id string) (RecordEntity, error)

	Delete(ctx context.Context, id string) error

	Find(ctx context.Context, q dal.Query) (QueryResult, error)

	FindByScriptID(ctx context.Context, scriptID string, q dal.Query) (QueryResult, error)
}
