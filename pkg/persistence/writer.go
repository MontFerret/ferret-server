package persistence

import (
	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/rs/zerolog"
)

type StdoutWriterFn func(job execution.Job, data []byte) error

func (fn StdoutWriterFn) Write(job execution.Job, data []byte) error {
	return fn(job, data)
}

func NewStdoutWriterFn(logger *zerolog.Logger) execution.OutputWriter {
	return StdoutWriterFn(func(job execution.Job, data []byte) error {
		logger.Info().
			Timestamp().
			Str("project_id", job.ProjectID).
			Str("job_id", job.ID).
			Bytes("data", data).
			Msg("job output")

		return nil
	})
}
