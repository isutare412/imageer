package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	"github.com/isutare412/imageer/internal/gateway/port"
)

type LeaderElector struct {
	clientSet *kubernetes.Clientset
	handlers  []port.LeaderHandler
	cfg       LeaderElectorConfig
	hostname  string

	wg             *sync.WaitGroup
	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
}

func NewLeaderElector(cfg LeaderElectorConfig, client *Client, handlers []port.LeaderHandler,
) (*LeaderElector, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("getting hostname: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &LeaderElector{
		clientSet:      client.clientSet,
		handlers:       handlers,
		cfg:            cfg,
		hostname:       hostname,
		wg:             &sync.WaitGroup{},
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
	}, nil
}

func (e *LeaderElector) Run() {
	lock := &resourcelock.LeaseLock{
		Client: e.clientSet.CoordinationV1(),
		LeaseMeta: metav1.ObjectMeta{
			Name:      e.cfg.LeaseName,
			Namespace: e.cfg.LeaseNamespace,
		},
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: e.hostname,
		},
	}

	go func() {
		for {
			select {
			case <-e.lifetimeCtx.Done():
				return
			default:
			}

			leaderelection.RunOrDie(e.lifetimeCtx, leaderelection.LeaderElectionConfig{
				Lock:          lock,
				LeaseDuration: e.cfg.LeaseDuration,
				RenewDeadline: e.cfg.RenewDeadline,
				RetryPeriod:   e.cfg.RetryPeriod,
				Callbacks: leaderelection.LeaderCallbacks{
					OnStartedLeading: e.onStartedLeading,
					OnStoppedLeading: e.onStoppedLeading,
				},
			})
		}
	}()
}

func (e *LeaderElector) Shutdown() {
	e.lifetimeCancel()
	e.wg.Wait()
}

func (e *LeaderElector) onStartedLeading(ctx context.Context) {
	slog.Info("Start leading as a leader")

	for _, h := range e.handlers {
		e.wg.Go(func() {
			h.OnStartedLeading(ctx)
		})
	}
}

func (e *LeaderElector) onStoppedLeading() {
	slog.Info("Stop leading as a leader")

	for _, h := range e.handlers {
		e.wg.Go(func() {
			h.OnStoppedLeading()
		})
	}
}
