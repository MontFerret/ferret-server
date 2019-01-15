package persistence

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret-server/pkg/execution"
)

type WriterFn func(job execution.Job, data []byte) error

func (fn WriterFn) Write(job execution.Job, data []byte) error {
	return fn(job, data)
}

func ToWriter(srv *Service) execution.OutputWriter {
	return WriterFn(func(job execution.Job, data []byte) error {
		if job.Script.Persistence.Enabled {
			var val interface{}

			err := json.Unmarshal(data, &val)

			if err != nil {
				return err
			}

			_, err = srv.CreateRecord(context.Background(), job, val)

			return err
		}

		return nil
	})
}
