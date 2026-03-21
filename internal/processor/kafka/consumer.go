package kafka

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/isutare412/imageer/pkg/kafkahelpers"
	"github.com/isutare412/imageer/pkg/log"
)

const headerKeyRetryCount = "retry-count"

type Consumer struct {
	client   *kgo.Client
	handlers map[string]Handler

	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
	wg             sync.WaitGroup
}

func NewConsumer(client *Client, handlers map[string]Handler) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = log.WithAttrContext(ctx)

	return &Consumer{
		client:         client.inner,
		handlers:       handlers,
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
	}
}

func (c *Consumer) Run() {
	c.wg.Go(c.pollLoop)
}

func (c *Consumer) Shutdown() {
	c.lifetimeCancel()
	c.wg.Wait()
}

func (c *Consumer) pollLoop() {
	for {
		fetches := c.client.PollFetches(c.lifetimeCtx)
		if c.lifetimeCtx.Err() != nil {
			slog.Info("Break kafka consumer poll loop")
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				slog.Error("Kafka fetch error",
					"topic", err.Topic, "partition", err.Partition, "error", err.Err)
			}
		}

		// NOTE: We fan out a goroutine per partition to process records
		// concurrently across partitions while preserving ordering within
		// each partition.
		var wg sync.WaitGroup
		fetches.EachPartition(func(ftp kgo.FetchTopicPartition) {
			handler, ok := c.handlers[ftp.Topic]
			if !ok {
				slog.Warn("No handler for topic", "topic", ftp.Topic)
				return
			}

			wg.Go(func() {
				for _, record := range ftp.Records {
					handler.HandleRecord(context.WithoutCancel(c.lifetimeCtx), record)
				}
			})
		})
		wg.Wait()
	}
}

func (c *Consumer) scheduleRetry(handler Handler, record *kgo.Record, retryCount int) {
	c.wg.Go(func() {
		delay := handler.RetryBaseDelay() * time.Duration(retryCount)
		select {
		case <-time.After(delay):
		case <-c.lifetimeCtx.Done():
		}

		retryRecord := &kgo.Record{
			Topic: handler.RetryTopic(),
			Value: record.Value,
			Headers: []kgo.RecordHeader{
				{Key: headerKeyRetryCount, Value: []byte(strconv.Itoa(retryCount))},
			},
		}

		results := c.client.ProduceSync(context.Background(), retryRecord)
		if err := results.FirstErr(); err != nil {
			slog.Error("Failed to produce retry record",
				"topic", handler.RetryTopic(), "retryCount", retryCount,
				"error", kafkahelpers.WrapKafkaError(err, "Failed to produce retry"))
		}

		slog.Info("Produced kafka retry record",
			"topic", handler.RetryTopic(), "retryCount", retryCount)
	})
}

func parseRetryCount(record *kgo.Record) int {
	for _, h := range record.Headers {
		if h.Key == headerKeyRetryCount {
			count, err := strconv.Atoi(string(h.Value))
			if err != nil {
				return 0
			}
			return count
		}
	}
	return 0
}
