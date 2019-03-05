package execution_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/MontFerret/ferret-server/pkg/execution"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/smartystreets/goconvey/convey"
)

type (
	NoopStatusWriter struct{}

	MockedStatusWriter struct {
		mu   sync.Mutex
		jobs map[string][]execution.Status
	}

	MockedWorker struct {
		mu        sync.Mutex
		job       execution.Job
		data      func() []byte
		err       func() error
		running   bool
		cancelled bool
	}
)

func (nsw NoopStatusWriter) Write(state execution.State) error {
	return nil
}

func (nw *MockedWorker) IsRunning() bool {
	nw.mu.Lock()
	defer nw.mu.Unlock()

	return nw.running
}

func (nw *MockedWorker) Process() ([]byte, error) {
	nw.mu.Lock()
	nw.running = true
	nw.mu.Unlock()

	time.Sleep(time.Duration(100) * time.Millisecond)

	nw.mu.Lock()
	nw.running = false
	defer nw.mu.Unlock()

	var data []byte
	var err error

	if nw.data != nil {
		data = nw.data()
	}

	if nw.err != nil {
		err = nw.err()
	}

	return data, err
}

func (nw *MockedWorker) Interrupt() {
	nw.mu.Lock()
	defer nw.mu.Unlock()

	nw.cancelled = true
}

func (sw *MockedStatusWriter) Write(state execution.State) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	found, ok := sw.jobs[state.Job.ID]
	if !ok {
		found = make([]execution.Status, 0, 5)
	}

	sw.jobs[state.Job.ID] = append(found, state.Status)

	return nil
}

func NewNoopLogger() *zerolog.Logger {
	l := zerolog.New(NoopWriter{})

	return &l
}

func NewNoopStatusWriter() execution.StateWriter {
	return NoopStatusWriter{}
}

func NewMockedWorker(job execution.Job, data func() []byte, err func() error) execution.Worker {
	w := new(MockedWorker)
	w.job = job
	w.data = data
	w.err = err

	return w
}

func NewMockedStatusWriter() *MockedStatusWriter {
	sw := new(MockedStatusWriter)
	sw.jobs = make(map[string][]execution.Status)

	return sw
}

func TestPool(t *testing.T) {
	Convey(".Consume", t, func() {
		Convey("Should read messages from a given queue", func() {
			jobs := []execution.Job{
				MockNoopJob(),
				MockNoopJob(),
				MockNoopJob(),
				MockNoopJob(),
				MockNoopJob(),
				MockNoopJob(),
				MockNoopJob(),
				MockNoopJob(),
				MockNoopJob(),
				MockNoopJob(),
			}

			q, err := execution.NewInMemoryQueue(10)

			So(err, ShouldBeNil)

			p, err := execution.NewWorkerPool(
				1,
				NewNoopLogger(),
				NewNoopStatusWriter(),
				func(job execution.Job) execution.Worker {
					return NewMockedWorker(job, func() []byte {
						return []byte("foo")
					}, nil)
				},
			)

			So(err, ShouldBeNil)

			ctx := context.Background()

			res, err := p.Consume(ctx, q)

			So(err, ShouldBeNil)

			for _, in := range jobs {
				So(q.Enqueue(ctx, in), ShouldBeNil)
			}

			i := -1
			stopAt := len(jobs) - 1

			for i != stopAt {
				select {
				case j := <-res:
					i++
					So(j.Job.ID, ShouldEqual, jobs[i].ID)
				default:
					return
				}
			}
		})

		Convey("Should update job state", func() {
			completedJob := MockNoopJob()
			failedJob := MockNoopJob()
			jobs := []execution.Job{
				completedJob,
				failedJob,
			}

			q, err := execution.NewInMemoryQueue(10)

			So(err, ShouldBeNil)

			sw := NewMockedStatusWriter()

			p, err := execution.NewWorkerPool(
				1,
				NewNoopLogger(),
				sw,
				func(job execution.Job) execution.Worker {
					if job.ID == failedJob.ID {
						return NewMockedWorker(job, nil, func() error {
							return errors.New("test")
						})
					}

					return NewMockedWorker(job, func() []byte {
						return []byte("foo")
					}, nil)
				},
			)

			So(err, ShouldBeNil)

			ctx := context.Background()

			res, err := p.Consume(ctx, q)

			So(err, ShouldBeNil)

			for _, in := range jobs {
				So(q.Enqueue(ctx, in), ShouldBeNil)
			}

			i := -1
			stopAt := len(jobs) - 1

			for stopAt != i {
				select {
				case j := <-res:
					i++
					So(j.Job.ID, ShouldEqual, jobs[i].ID)
				default:
					break
				}
			}

			time.Sleep(time.Duration(1000) * time.Millisecond)

			sw.mu.Lock()
			completedStatuses := sw.jobs[completedJob.ID]
			sw.mu.Unlock()

			So(completedStatuses[0], ShouldEqual, execution.StatusRunning)
			So(completedStatuses[1], ShouldEqual, execution.StatusCompleted)

			sw.mu.Lock()
			failedStatuses := sw.jobs[failedJob.ID]
			sw.mu.Unlock()

			So(failedStatuses[0], ShouldEqual, execution.StatusRunning)
			So(failedStatuses[1], ShouldEqual, execution.StatusErrored)
		})

		Convey("Should cancel a running job by its ID", func() {
			q, err := execution.NewInMemoryQueue(10)

			So(err, ShouldBeNil)

			sw := NewMockedStatusWriter()

			c := make(chan []byte)

			p, err := execution.NewWorkerPool(
				1,
				NewNoopLogger(),
				sw,
				func(job execution.Job) execution.Worker {
					return NewMockedWorker(job, func() []byte {
						out := <-c

						return out
					}, nil)
				},
			)

			So(err, ShouldBeNil)

			ctx := context.Background()

			_, err = p.Consume(ctx, q)

			So(err, ShouldBeNil)

			j := MockJob(`RETURN TRUE`, make(map[string]interface{}))

			q.Enqueue(ctx, j)

			time.Sleep(time.Duration(100) * time.Millisecond)

			So(p.Cancel(j.ProjectID, j.ID), ShouldBeNil)

			c <- []byte("foo")

			time.Sleep(time.Duration(100) * time.Millisecond)

			sw.mu.Lock()
			statuses := sw.jobs[j.ID]
			sw.mu.Unlock()

			So(statuses[0], ShouldEqual, execution.StatusRunning)
			So(statuses[1], ShouldEqual, execution.StatusCancelled)
		})
	})
}
