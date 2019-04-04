package scripts

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/pkg/errors"
)

type (
	DbContext interface {
		GetScriptsRepository(projectID string) (Repository, error)
	}

	Service struct {
		db DbContext
	}
)

func NewService(db DbContext) (*Service, error) {
	if db == nil {
		return nil, common.Error(common.ErrMissedArgument, "db context")
	}

	return &Service{db}, nil
}

func (service *Service) GetScript(ctx context.Context, projectID, id string) (ScriptEntity, error) {
	repo, err := service.db.GetScriptsRepository(projectID)

	if err != nil {
		return ScriptEntity{}, errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "scripts")
	}

	return repo.Get(ctx, id)
}

func (service *Service) FindScripts(ctx context.Context, projectID string, q dal.Query) (QueryResult, error) {
	repo, err := service.db.GetScriptsRepository(projectID)

	if err != nil {
		return QueryResult{}, errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "scripts")
	}

	return repo.Find(ctx, q)
}

func (service *Service) CreateScript(ctx context.Context, projectID string, script Script) (dal.Entity, error) {
	repo, err := service.db.GetScriptsRepository(projectID)

	if err != nil {
		return dal.Entity{}, errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "scripts")
	}

	return repo.Create(ctx, script)
}

func (service *Service) UpdateScript(ctx context.Context, projectID string, script UpdateScript) (dal.Entity, error) {
	repo, err := service.db.GetScriptsRepository(projectID)

	if err != nil {
		return dal.Entity{}, errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "scripts")
	}

	return repo.Update(ctx, script)
}

func (service *Service) DeleteScript(ctx context.Context, projectID, id string) error {
	repo, err := service.db.GetScriptsRepository(projectID)

	if err != nil {
		return errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "scripts")
	}

	return repo.Delete(ctx, id)
}
