package history

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
)

type Repository interface {
	Create(ctx context.Context, entry Record) (dal.Entity, error)

	Update(ctx context.Context, entry Record) (dal.Entity, error)

	Get(ctx context.Context, id string) (RecordEntity, error)

	Find(ctx context.Context, q dal.Query) (QueryResult, error)
}
