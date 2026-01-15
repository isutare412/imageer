package kubernetes

import (
	"context"
	"log/slog"
	"sync"

	"github.com/isutare412/imageer/internal/gateway/port"
)

// StandaloneElector mimics LeaderElector behavior for development environments
// without Kubernetes. It immediately acts as a leader without election.
type StandaloneElector struct {
	handlers []port.LeaderHandler

	wg             *sync.WaitGroup
	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
}

func NewStandaloneElector(handlers []port.LeaderHandler) *StandaloneElector {
	ctx, cancel := context.WithCancel(context.Background())

	return &StandaloneElector{
		handlers:       handlers,
		wg:             &sync.WaitGroup{},
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
	}
}

func (e *StandaloneElector) Run() {
	e.onStartedLeading(e.lifetimeCtx)
}

func (e *StandaloneElector) Shutdown() {
	e.lifetimeCancel()
	e.onStoppedLeading()
	e.wg.Wait()
}

func (e *StandaloneElector) onStartedLeading(ctx context.Context) {
	slog.Info("Start leading as a leader")

	for _, h := range e.handlers {
		e.wg.Go(func() {
			h.OnStartedLeading(ctx)
		})
	}
}

func (e *StandaloneElector) onStoppedLeading() {
	slog.Info("Stop leading as a leader")

	for _, h := range e.handlers {
		e.wg.Go(func() {
			h.OnStoppedLeading()
		})
	}
}
