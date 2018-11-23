package execution_test

import (
	"context"
	"github.com/MontFerret/ferret-server/pkg/execution"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInMemoryQueue(t *testing.T) {
	Convey("Should pass messages", t, func() {
		ctx := context.Background()
		q, err := execution.NewInMemoryQueue(1)

		So(err, ShouldBeNil)

		job := MockNoopJob()
		err = q.Enqueue(ctx, job)
		So(err, ShouldBeNil)

		ch, err := q.Dequeue(ctx)

		So(err, ShouldBeNil)
		So(ch, ShouldNotBeNil)

		outJ := <-ch

		So(job.ID, ShouldEqual, outJ.ID)
	})

	Convey("Should return an error when queue is full", t, func() {
		ctx := context.Background()
		q, err := execution.NewInMemoryQueue(1)

		So(err, ShouldBeNil)

		err = q.Enqueue(ctx, MockNoopJob())
		So(err, ShouldBeNil)

		err = q.Enqueue(ctx, MockNoopJob())
		So(err, ShouldNotBeNil)

		ch, err := q.Dequeue(context.Background())

		So(err, ShouldBeNil)
		So(ch, ShouldNotBeNil)
	})

	Convey("Should pass messages in FIFO order", t, func() {
		ctx := context.Background()
		jobs := []execution.Job{
			MockNoopJob(),
			MockNoopJob(),
			MockNoopJob(),
			MockNoopJob(),
			MockNoopJob(),
		}

		q, err := execution.NewInMemoryQueue(10)
		So(err, ShouldBeNil)

		for _, j := range jobs {
			So(q.Enqueue(ctx, j), ShouldBeNil)
		}

		ch, err := q.Dequeue(ctx)
		So(err, ShouldBeNil)

		i := -1

		for {
			select {
			case j := <-ch:
				i++
				So(j.ID, ShouldEqual, jobs[i].ID)
			default:
				return
			}
		}
	})
}
