package projects

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/pkg/errors"
)

type (
	DbContext interface {
		GetProjectsRepository() (Repository, error)
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

func (service *Service) GetProject(ctx context.Context, id string) (ProjectEntity, error) {
	repo, err := service.db.GetProjectsRepository()

	if err != nil {
		return ProjectEntity{}, errors.Wrap(err, "resolve project repository")
	}

	return repo.Get(ctx, id)
}

func (service *Service) FindProjects(ctx context.Context, q dal.Query) (QueryResult, error) {
	repo, err := service.db.GetProjectsRepository()

	if err != nil {
		return QueryResult{}, errors.Wrap(err, "resolve project repository")
	}

	return repo.Find(ctx, q)
}

func (service *Service) CreateProject(ctx context.Context, project Project) (dal.Entity, error) {
	repo, err := service.db.GetProjectsRepository()

	if err != nil {
		return dal.Entity{}, errors.Wrap(err, "resolve project repository")
	}

	return repo.Create(ctx, project)
}

func (service *Service) UpdateProject(ctx context.Context, project UpdateProject) (dal.Entity, error) {
	repo, err := service.db.GetProjectsRepository()

	if err != nil {
		return dal.Entity{}, errors.Wrap(err, "resolve project repository")
	}

	return repo.Update(ctx, project)
}

func (service *Service) DeleteProject(ctx context.Context, id string) error {
	repo, err := service.db.GetProjectsRepository()

	if err != nil {
		return errors.Wrap(err, "resolve project repository")
	}

	return repo.Delete(ctx, id)
}
