package execution

import (
	"context"
	"runtime"
	"time"

	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type (
	DbContext interface {
		GetScriptsRepository(projectID string) (scripts.Repository, error)
	}

	StateWriter interface {
		Write(state State) error
	}

	LogWriter interface {
		Write(job Job, data []byte) (n int, err error)
	}

	OutputWriter interface {
		Write(job Job, data []byte) error
	}

	Service struct {
		logger   *zerolog.Logger
		db       DbContext
		compiler *compiler.FqlCompiler
		queue    Queue
		pool     *WorkerPool
		state    StateWriter
		output   OutputWriter
	}
)

func NewService(
	settings Settings,
	logger *zerolog.Logger,
	db DbContext,
	compiler *compiler.FqlCompiler,
	queue Queue,
	state StateWriter,
	logs LogWriter,
	output OutputWriter,
) (*Service, error) {
	if logger == nil {
		return nil, common.Error(common.ErrMissedArgument, "logger")
	}

	if db == nil {
		return nil, common.Error(common.ErrMissedArgument, "db context")
	}

	if compiler == nil {
		return nil, common.Error(common.ErrMissedArgument, "FQL compiler")
	}

	if state == nil {
		return nil, common.Error(common.ErrMissedArgument, "state writer")
	}

	if output == nil {
		return nil, common.Error(common.ErrMissedArgument, "output writer")
	}

	s := new(Service)
	s.logger = logger
	s.db = db
	s.compiler = compiler
	s.queue = queue
	s.state = state
	s.output = output

	size := uint64(runtime.NumCPU() * int(settings.PoolSize))
	pool, err := NewWorkerPool(size, logger, state, func(job Job) Worker {
		return NewFQLWorker(compiler, NewLogger(logs, job), job)
	})

	if err != nil {
		return nil, err
	}

	s.pool = pool

	out, err := s.pool.Consume(context.Background(), queue)

	if err != nil {
		return nil, err
	}

	s.start(out)

	return s, nil
}

func (service *Service) Start(ctx context.Context, event Event) (string, error) {
	scriptsRepo, err := service.db.GetScriptsRepository(event.ProjectID)

	if err != nil {
		return "", errors.Wrapf(err, "%s %s", dal.ErrResolveRepo, "scripts")
	}

	entity, err := scriptsRepo.Get(ctx, event.ScriptID)

	if err != nil {
		return "", errors.Wrapf(err, "get entity %s", event.ScriptID)
	}

	id, err := uuid.NewV4()

	if err != nil {
		return "", errors.Wrap(err, "create an jobID")
	}

	jobID := id.String()

	job := Job{
		ID:        jobID,
		ProjectID: event.ProjectID,
		Script:    entity,
		CausedBy:  event.CausedBy,
	}

	err = service.state.Write(State{
		Job:       job,
		Status:    StatusQueued,
		Timestamp: time.Now(),
	})

	if err != nil {
		service.log(job, err).
			Str("state", StatusQueued.String()).
			Msg("failed to create job state")
	}

	err = service.queue.Enqueue(ctx, job)

	if err != nil {
		e := service.state.Write(State{
			Job:       job,
			Status:    StatusErrored,
			Timestamp: time.Now(),
			Error:     err,
		})

		if e != nil {
			service.log(job, err).
				Str("state", StatusQueued.String()).
				Msg("failed to update job state")
		}

		return "", errors.Wrap(err, "enqueue")
	}

	return jobID, nil
}

func (service *Service) Cancel(_ context.Context, projectID string, jobID string) error {
	return service.pool.Cancel(projectID, jobID)
}

func (service *Service) start(out <-chan Result) {
	go func() {
		for result := range out {
			err := service.output.Write(result.Job, result.Data)

			if err != nil {
				service.log(result.Job, err).
					Msg("failed to save script result data")
			}
		}
	}()
}

func (service *Service) log(job Job, err error) *zerolog.Event {
	if err != nil {
		return service.logger.
			Error().
			Timestamp().
			Err(err).
			Str("project_id", job.ProjectID).
			Str("job_id", job.ID)
	}

	return service.logger.
		Info().
		Timestamp().
		Str("project_id", job.ProjectID).
		Str("job_id", job.ID)
}
