package persistence

import (
	"context"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/MontFerret/ferret-server/pkg/scripts"

	"github.com/pkg/errors"
)

type (
	DbContext interface {
		GetScriptsRepository(projectID string) (scripts.Repository, error)
		GetDataRepository(projectID string) (Repository, error)
	}

	Service struct {
		db DbContext
	}
)

func NewService(db DbContext) (*Service, error) {
	if db == nil {
		return nil, common.Error(common.ErrMissedArgument, "db")
	}

	return &Service{db: db}, nil
}

func (service *Service) CreateRecord(ctx context.Context, job execution.Job, data interface{}) (dal.Entity, error) {
	if job.Script.Persistence.Enabled == false {
		return dal.Entity{}, errors.Errorf(
			"persistence is disabled for a given script: %s rev %s",
			job.Script.ID,
			job.Script.Rev,
		)
	}

	repo, err := service.resolveRepository(ctx, job.ProjectID)

	if err != nil {
		return dal.Entity{}, err
	}

	entity, err := repo.Create(ctx, Record{
		JobID:     job.ID,
		ScriptID:  job.Script.ID,
		ScriptRev: job.Script.Rev,
		Data:      data,
	})

	if err != nil {
		return dal.Entity{}, errors.Wrap(err, "create a new record")
	}

	return entity, err
}

func (service *Service) UpdateRecord(ctx context.Context, projectID string, record Record) (dal.Entity, error) {
	repo, err := service.resolveRepository(ctx, projectID)

	if err != nil {
		return dal.Entity{}, err
	}

	entity, err := repo.Update(ctx, record)

	if err != nil {
		return dal.Entity{}, errors.Wrap(err, "update a record")
	}

	return entity, nil
}

func (service *Service) DeleteRecord(ctx context.Context, identity scripts.Identity, id string) error {
	repo, err := service.resolveRepository(ctx, identity.ProjectID)

	if err != nil {
		return err
	}

	err = repo.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (service *Service) GetRecord(ctx context.Context, identity scripts.Identity, id string) (RecordEntity, error) {
	repo, err := service.resolveRepository(ctx, identity.ProjectID)

	if err != nil {
		return RecordEntity{}, err
	}

	return repo.Get(ctx, id)
}

func (service *Service) FindScriptRecords(ctx context.Context, identity scripts.Identity, q dal.Query) (QueryResult, error) {
	repo, err := service.resolveRepository(ctx, identity.ProjectID)

	if err != nil {
		return QueryResult{}, err
	}

	return repo.FindByScriptID(ctx, identity.ScriptID, q)
}

func (service *Service) FindProjectRecords(ctx context.Context, projectID string, q dal.Query) (QueryResult, error) {
	repo, err := service.resolveRepository(ctx, projectID)

	if err != nil {
		return QueryResult{}, err
	}

	out, err := repo.Find(ctx, q)

	if err != nil {
		return QueryResult{}, err
	}

	return out, nil
}

func (service *Service) resolveRepository(_ context.Context, projectID string) (Repository, error) {
	repo, err := service.db.GetDataRepository(projectID)

	if err != nil {
		return nil, errors.Wrap(err, "resolve a repository")
	}

	return repo, nil
}
