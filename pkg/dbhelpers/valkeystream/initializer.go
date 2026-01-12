package valkeystream

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

// Initializer initializes valkey consumer group and consumer before use.
type Initializer struct {
	client valkey.Client
	cfg    InitializerConfig
}

func NewInitializer(client valkey.Client, cfg InitializerConfig) *Initializer {
	return &Initializer{
		client: client,
		cfg:    cfg,
	}
}

func (i *Initializer) Initialize(ctx context.Context) error {
	// Create consumer group if not exists
	resp := i.client.Do(ctx, i.client.B().XgroupCreate().
		Key(i.cfg.Consumer.Stream).
		Group(i.cfg.Consumer.Group).
		Id("0").
		Mkstream().
		Build())
	err := dbhelpers.WrapValkeyError(resp.Error(), "Failed to create consumer group %s",
		i.cfg.Consumer.Group)
	switch {
	case apperr.IsErrorStatusCode(err, http.StatusConflict):
		// Skip if group already exists
	case err != nil:
		return err
	}

	// Create consumer in the group
	resp = i.client.Do(ctx, i.client.B().XgroupCreateconsumer().
		Key(i.cfg.Consumer.Stream).
		Group(i.cfg.Consumer.Group).
		Consumer(i.cfg.Consumer.Name).
		Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to create consumer %s in group %s",
			i.cfg.Consumer.Name, i.cfg.Consumer.Group)
	}

	slog.InfoContext(ctx, "Created valkey consumer",
		"consumer", i.cfg.Consumer.Name, "consumerGroup", i.cfg.Consumer.Group,
		"stream", i.cfg.Consumer.Stream)
	return nil
}
