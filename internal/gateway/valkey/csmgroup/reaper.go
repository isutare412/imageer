package csmgroup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type Reaper struct {
	client            valkey.Client
	stream            string
	group             string
	idleTimeThreshold time.Duration
}

func NewReaper(client valkey.Client, stream, group string, idleTimeThreshold time.Duration,
) *Reaper {
	return &Reaper{
		client:            client,
		stream:            stream,
		group:             group,
		idleTimeThreshold: idleTimeThreshold,
	}
}

func (r *Reaper) ReapIdleConsumers(ctx context.Context) error {
	// List consumers in the group
	resp := r.client.Do(ctx, r.client.B().XinfoConsumers().
		Key(r.stream).
		Group(r.group).
		Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to get consumers info for group %s", r.group)
	}

	consumersToReap, err := findConsumersToReap(resp, r.idleTimeThreshold)
	if err != nil {
		return fmt.Errorf("finding consumers to reap: %w", err)
	}

	// Reap dead consumers
	for _, name := range consumersToReap {
		resp := r.client.Do(ctx, r.client.B().XgroupDelconsumer().
			Key(r.stream).
			Group(r.group).
			Consumername(name).
			Build())
		if err := resp.Error(); err != nil {
			return dbhelpers.WrapValkeyError(err, "Failed to delete consumer %s from group %s", name, r.group)
		}

		slog.InfoContext(ctx, "Reaped idle valkey consumer",
			"consumer", name, "consumerGroup", r.group, "stream", r.stream)
	}
	return nil
}
