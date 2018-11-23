package history

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/execution"
)

type StatusWriterFn func(state execution.State) error

func (fn StatusWriterFn) Write(status execution.State) error {
	return fn(status)
}

func ToStatusWriter(service *Service) execution.StateWriter {
	return StatusWriterFn(func(state execution.State) error {
		if state.Status == execution.StatusQueued {
			return service.Create(context.Background(), state.Job)
		}

		_, err := service.Update(context.Background(), state)

		return err
	})
}

type LogWriterFn func(job execution.Job, data []byte) (n int, err error)

func (fn LogWriterFn) Write(job execution.Job, data []byte) (n int, err error) {
	return fn(job, data)
}

func ToLogWriter(service *Service) execution.LogWriter {
	return LogWriterFn(func(job execution.Job, data []byte) (n int, err error) {
		_, err = service.Log(context.Background(), job.ProjectID, job.ID, data)

		return 1, err
	})
}
