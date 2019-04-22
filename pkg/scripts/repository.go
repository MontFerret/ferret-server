package scripts

import (
	"context"

	"github.com/MontFerret/ferret-server/pkg/common/dal"
)

type Repository interface {
	Get(ctx context.Context, id string) (ScriptEntity, error)

	Find(ctx context.Context, query dal.Query) (QueryResult, error)

	Create(ctx context.Context, script Script) (dal.Entity, error)

	Update(ctx context.Context, script UpdateScript) (dal.Entity, error)

	Delete(ctx context.Context, id string) error
}
