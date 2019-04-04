package history

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/pkg/errors"
)

type (
	DbContext interface {
		GetHistoryRepository(projectID string) (Repository, error)
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

func (service *Service) Create(ctx context.Context, job execution.Job) error {
	repo, err := service.db.GetHistoryRepository(job.ProjectID)

	if err != nil {
		return errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "history")
	}

	record := Record{}
	record.JobID = job.ID
	record.ScriptID = job.Script.ID
	record.ScriptRev = job.Script.Rev
	record.Params = job.Script.Execution.Params
	record.Status = execution.StatusQueued
	record.CausedBy = job.CausedBy

	_, err = repo.Create(ctx, record)

	return err
}

func (service *Service) Log(ctx context.Context, projectID, jobID string, data []byte) (dal.Entity, error) {
	repo, err := service.db.GetHistoryRepository(projectID)

	if err != nil {
		return dal.Entity{}, errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "history")
	}

	found, err := repo.Get(ctx, jobID)

	if err != nil {
		return dal.Entity{}, err
	}

	if found.Logs == nil {
		found.Logs = make([]string, 0, 10)
	}

	found.Logs = append(found.Logs, string(data))

	return repo.Update(ctx, found.Record)
}

func (service *Service) Update(ctx context.Context, state execution.State) (dal.Entity, error) {
	repo, err := service.db.GetHistoryRepository(state.Job.ProjectID)

	if err != nil {
		return dal.Entity{}, errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "history")
	}

	found, err := repo.Get(ctx, state.Job.ID)

	if err != nil {
		return dal.Entity{}, err
	}

	found.Status = state.Status

	switch state.Status {
	case execution.StatusRunning:
		found.StartedAt = state.Timestamp
	case execution.StatusCompleted, execution.StatusCancelled:
		found.EndedAt = state.Timestamp
	case execution.StatusErrored:
		found.EndedAt = state.Timestamp
		found.Error = state.Error
	}

	return repo.Update(ctx, found.Record)
}

func (service *Service) Get(ctx context.Context, projectID, jobID string) (RecordEntity, error) {
	repo, err := service.db.GetHistoryRepository(projectID)

	if err != nil {
		return RecordEntity{}, errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "history")
	}

	return repo.Get(ctx, jobID)
}

func (service *Service) Find(ctx context.Context, projectID string, q dal.Query) (QueryResult, error) {
	repo, err := service.db.GetHistoryRepository(projectID)

	if err != nil {
		return QueryResult{}, errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "history")
	}

	return repo.Find(ctx, q)
}
