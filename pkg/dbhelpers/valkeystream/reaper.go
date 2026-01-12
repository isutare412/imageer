package valkeystream

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/dbhelpers"
)

// Reaper reaps idle consumers from a consumer group.
type Reaper struct {
	client valkey.Client
	cfg    ReaperConfig
}

func NewReaper(client valkey.Client, cfg ReaperConfig) *Reaper {
	return &Reaper{
		client: client,
		cfg:    cfg,
	}
}

func (r *Reaper) ReapIdleConsumers(ctx context.Context) error {
	// List consumers in the group
	resp := r.client.Do(ctx, r.client.B().XinfoConsumers().
		Key(r.cfg.Stream).
		Group(r.cfg.Group).
		Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to get consumers info for group %s",
			r.cfg.Group)
	}

	consumersToReap, err := findConsumersToReap(resp, r.cfg.IdleTimeThreshold)
	if err != nil {
		return fmt.Errorf("finding consumers to reap: %w", err)
	}

	// Reap dead consumers
	for _, name := range consumersToReap {
		resp := r.client.Do(ctx, r.client.B().XgroupDelconsumer().
			Key(r.cfg.Stream).
			Group(r.cfg.Group).
			Consumername(name).
			Build())
		if err := resp.Error(); err != nil {
			return dbhelpers.WrapValkeyError(err, "Failed to delete consumer %s from group %s",
				name, r.cfg.Group)
		}

		slog.InfoContext(ctx, "Reaped idle valkey consumer",
			"consumer", name, "consumerGroup", r.cfg.Group, "stream", r.cfg.Stream)
	}
	return nil
}
