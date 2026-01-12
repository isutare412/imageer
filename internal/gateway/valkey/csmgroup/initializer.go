package csmgroup

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type Initializer struct {
	client   valkey.Client
	stream   string
	group    string
	consumer string
}

func NewInitializer(client valkey.Client, stream, group, consumer string) *Initializer {
	return &Initializer{
		client:   client,
		stream:   stream,
		group:    group,
		consumer: consumer,
	}
}

func (i *Initializer) Initialize(ctx context.Context) error {
	// Create consumer group if not exists
	resp := i.client.Do(ctx, i.client.B().XgroupCreate().
		Key(i.stream).
		Group(i.group).
		Id("0").
		Mkstream().
		Build())
	err := dbhelpers.WrapValkeyError(resp.Error(), "Failed to create consumer group %s", i.group)
	switch {
	case apperr.IsErrorStatusCode(err, http.StatusConflict):
		// Skip if group already exists
	case err != nil:
		return err
	}

	// Create consumer in the group
	resp = i.client.Do(ctx, i.client.B().XgroupCreateconsumer().
		Key(i.stream).
		Group(i.group).
		Consumer(i.consumer).
		Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to create consumer %s in group %s", i.consumer, i.group)
	}

	slog.InfoContext(ctx, "Created valkey consumer",
		"consumer", i.consumer, "consumerGroup", i.group, "stream", i.stream)

	return nil
}
