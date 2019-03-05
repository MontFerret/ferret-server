package execution_test

import (
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/MontFerret/ferret-server/pkg/common/dal"
	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"github.com/gofrs/uuid"

	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
)

type NoopWriter struct{}

func (l NoopWriter) Write(data []byte) (int, error) {
	return 0, nil
}

func NewNoopWriter() io.Writer {
	return NoopWriter{}
}

func MockJob(q string, params map[string]interface{}) execution.Job {
	jobID, _ := uuid.NewV4()
	projectID, _ := uuid.NewV4()

	return execution.Job{
		ID:        jobID.String(),
		ProjectID: projectID.String(),
		CausedBy:  execution.CauseManual,
		Script:    MockScript(q, params),
	}
}

func MockNoopJob() execution.Job {
	return MockJob(`RETURN FALSE`, make(map[string]interface{}))
}

func MockScript(query string, params map[string]interface{}) scripts.ScriptEntity {
	scriptID, _ := uuid.NewV4()
	rev, _ := uuid.NewV4()

	return scripts.ScriptEntity{
		Entity: dal.Entity{
			ID:  scriptID.String(),
			Rev: rev.String(),
			Metadata: dal.Metadata{
				CreatedAt: time.Now(),
				UpdateAt:  time.Now(),
			},
		},
		Script: scripts.Script{
			Name:        "mock_script",
			Description: "mock script",
			Execution: scripts.Execution{
				Query:  query,
				Params: params,
			},
		},
	}
}

func TestWorker(t *testing.T) {
	c := compiler.New()

	createWorker := func(q string, params map[string]interface{}) execution.Worker {
		return execution.NewFQLWorker(c, NewNoopWriter(), MockJob(q, params))
	}

	Convey(".Process", t, func() {
		Convey("Should run worker", func() {
			w := createWorker(`
					FOR i IN 1..@size
						RETURN i
				`, map[string]interface{}{
				"size": 1000,
			})

			out, err := w.Process()

			So(err, ShouldBeNil)
			So(out, ShouldNotBeNil)

			arr := make([]int, 0, 1000)

			json.Unmarshal(out, &arr)

			So(arr, ShouldHaveLength, 1000)
		})

		Convey("Should handle error", func() {
			w := createWorker(`
					FOR i IN 1..@size
						RETURN foo
				`, map[string]interface{}{
				"size": 1000,
			})

			out, err := w.Process()

			So(err, ShouldNotBeNil)
			So(out, ShouldBeNil)
		})
	})

	Convey(".IsRunning", t, func() {
		Convey("Should return true when a worker is running", func() {
			w := createWorker(`
					WAIT(1000)
					RETURN TRUE
				`, map[string]interface{}{
				"size": 1000,
			})

			go func() {
				w.Process()
			}()

			time.Sleep(time.Duration(100) * time.Millisecond)

			So(w.IsRunning(), ShouldBeTrue)
		})

		Convey("Should return false when a worker is not running", func() {
			w := createWorker(`
					RETURN TRUE
				`, map[string]interface{}{
				"size": 1000,
			})

			So(w.IsRunning(), ShouldBeFalse)
			w.Process()
			So(w.IsRunning(), ShouldBeFalse)
		})
	})
}
