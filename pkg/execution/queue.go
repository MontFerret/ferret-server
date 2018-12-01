package execution

import (
	"context"
	"errors"
	"github.com/MontFerret/ferret-server/pkg/common"
)

type (
	Queue interface {
		Enqueue(ctx context.Context, job Job) error
		Dequeue(ctx context.Context) (<-chan Job, error)
	}

	InMemoryQueue struct {
		queue chan Job
	}
)

func NewInMemoryQueue(size uint64) (*InMemoryQueue, error) {
	if size == 0 {
		return nil, common.Error(common.ErrInvalidOperation, "queue size cannot be 0")
	}

	q := new(InMemoryQueue)
	q.queue = make(chan Job, size)

	return q, nil
}

func (q *InMemoryQueue) Enqueue(ctx context.Context, job Job) error {
	select {
	case <-ctx.Done():
		return common.ErrTerminated
	case q.queue <- job:
		return nil
	default:
		return errors.New("queue is full")
	}
}

func (q *InMemoryQueue) Dequeue(_ context.Context) (<-chan Job, error) {
	return q.queue, nil
}
