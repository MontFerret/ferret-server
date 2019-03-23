package execution

import (
	"context"
	"io"
	"sync"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/pkg/errors"
)

var ErrWorkerAlreadyRunning = errors.New("worker already running")

type FQLWorker struct {
	mu       sync.Mutex
	compiler *compiler.FqlCompiler
	log      io.Writer
	job      Job
	cancel   context.CancelFunc
}

func NewFQLWorker(compiler *compiler.FqlCompiler, log io.Writer, job Job) Worker {
	w := new(FQLWorker)
	w.compiler = compiler
	w.log = log
	w.job = job

	return w
}

func (w *FQLWorker) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.cancel != nil
}

func (w *FQLWorker) Process() ([]byte, error) {
	w.mu.Lock()

	if w.cancel != nil {
		w.mu.Unlock()
		return nil, ErrWorkerAlreadyRunning
	}
	ctx, cancelFn := context.WithCancel(context.Background())
	ctx = drivers.WithContext(ctx, http.NewDriver(), drivers.AsDefault())
	ctx = drivers.WithContext(ctx, cdp.NewDriver())
	w.cancel = cancelFn

	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		w.cancel = nil
		w.mu.Unlock()
	}()

	// TODO: Add caching for frequent scripts
	program, err := w.compiler.Compile(w.job.Script.Execution.Query)

	if err != nil {
		return nil, err
	}

	params := make(map[string]interface{}, len(w.job.Script.Execution.Params))

	w.mu.Lock()
	for k, v := range w.job.Script.Execution.Params {
		params[k] = v
	}
	w.mu.Unlock()

	out, err := program.Run(
		ctx,
		runtime.WithLog(w.log),
		runtime.WithParams(params),
	)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (w *FQLWorker) Interrupt() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.cancel != nil {
		w.cancel()
	}
}
