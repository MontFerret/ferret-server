package execution

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/common"
	"github.com/rs/zerolog"
	"time"
)

type (
	Worker interface {
		IsRunning() bool
		Process() ([]byte, error)
		Interrupt()
	}

	WorkerFactory func(job Job) Worker

	WorkerPool struct {
		logger  *zerolog.Logger
		factory WorkerFactory
		status  StateWriter
		size    uint64
		workers map[string]Worker
		pool    chan bool
		results chan Result
	}
)

func NewWorkerPool(
	size uint64,
	logger *zerolog.Logger,
	status StateWriter,
	factory WorkerFactory,
) (*WorkerPool, error) {
	if logger == nil {
		return nil, common.Error(common.ErrMissedArgument, "logger")
	}

	if status == nil {
		return nil, common.Error(common.ErrMissedArgument, "state writer")
	}

	if factory == nil {
		return nil, common.Error(common.ErrMissedArgument, "factory")
	}

	wp := new(WorkerPool)
	wp.logger = logger
	wp.factory = factory
	wp.status = status
	wp.size = size
	wp.workers = make(map[string]Worker)
	wp.pool = make(chan bool, size)
	wp.results = make(chan Result, size)

	return wp, nil
}

func (wp *WorkerPool) Consume(ctx context.Context, q Queue) (<-chan Result, error) {
	queue, err := q.Dequeue(ctx)

	if err != nil {
		return nil, err
	}

	out := make(chan Result, wp.size)

	go func() {
		for {
			select {
			case <-ctx.Done():
				// stop all running workers
				for _, w := range wp.workers {
					w.Interrupt()
				}

				close(out)

				return
			case job := <-queue:
				// acquiring a gorouting lock by sending a message to a channel
				// if channel is data of capacity i.e. the amount of available goroutings is 0
				// the operation will block until the capacity increases
				wp.pool <- true

				state := State{
					Job:       job,
					Timestamp: time.Now(),
					Status:    StatusRunning,
				}
				err := wp.status.Write(state)

				if err != nil {
					wp.logger.Error().
						Timestamp().
						Err(err).
						Str("project_id", job.ProjectID).
						Str("job_id", job.ID).
						Str("state", state.Status.String()).
						Msg("failed to update job state")
				}

				worker := wp.factory(job)
				wp.workers[job.ID] = worker

				// start a new gorouting with the worker
				go func() {
					// run the worker
					out, err := worker.Process()

					jr := Result{
						State: State{
							Job:       job,
							Timestamp: time.Now(),
							Status:    StatusCompleted,
							Error:     err,
						},
						Data: out,
					}

					if err != nil {
						jr.Status = StatusErrored
					}

					// sending data results
					wp.results <- jr

					// releasing the gorouting
					<-wp.pool
				}()
			case result := <-wp.results:
				worker, found := wp.workers[result.Job.ID]

				if !found {
					break
				}

				if worker.IsRunning() {
					worker.Interrupt()
				}

				delete(wp.workers, result.Job.ID)

				err := wp.status.Write(result.State)

				if err != nil {
					wp.logger.Error().
						Timestamp().
						Err(err).
						Str("project_id", result.Job.ProjectID).
						Str("job_id", result.Job.ID).
						Str("state", result.Status.String()).
						Msg("failed to update job state")
				}

				// in case of operation was terminated and the channel is closed
				select {
				case <-ctx.Done():
					return
				default:
					out <- result
					break
				}
			}
		}
	}()

	return out, nil
}

func (wp *WorkerPool) Cancel(projectID, jobID string) error {
	wp.results <- Result{
		State: State{
			Job: Job{
				ID:        jobID,
				ProjectID: projectID,
			},
			Status: StatusCancelled,
		},
	}

	return nil
}
