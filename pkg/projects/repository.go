package projects

import (
	"context"

	"github.com/MontFerret/ferret-server/pkg/common/dal"
)

type Repository interface {
	Get(ctx context.Context, id string) (ProjectEntity, error)

	Find(ctx context.Context, query dal.Query) (QueryResult, error)

	Create(ctx context.Context, project Project) (dal.Entity, error)

	Update(ctx context.Context, project UpdateProject) (dal.Entity, error)

	Delete(ctx context.Context, id string) error
}
