package csmgroup

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/apperr"
)

func GenerateConsumerName(prefix string) string {
	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)

	b32Enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	randString := b32Enc.EncodeToString(randBytes)
	return fmt.Sprintf("%s-consumer-%s", prefix, strings.ToLower(randString)[:10])
}

func findConsumersToReap(xinfoResult valkey.ValkeyResult, idleTimeThreshold time.Duration,
) ([]string, error) {
	consumers, err := xinfoResult.ToArray()
	if err != nil {
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to parse consumers list")
	}

	var names []string
	for _, consumer := range consumers {
		info, err := consumer.AsMap()
		if err != nil {
			return nil, apperr.NewError(apperr.CodeInternalServerError).
				WithCause(err).
				WithSummary("Failed to parse consumer map")
		}

		msg := info["pending"]
		pending, err := msg.AsInt64()
		if err != nil {
			return nil, apperr.NewError(apperr.CodeInternalServerError).
				WithCause(err).
				WithSummary("Failed to parse pending message count")
		}

		if pending > 0 {
			continue // Do not reap consumers with pending messages left
		}

		msg = info["idle"]
		idleMs, err := msg.AsInt64()
		if err != nil {
			return nil, apperr.NewError(apperr.CodeInternalServerError).
				WithCause(err).
				WithSummary("Failed to parse idle time")
		}

		idleTime := time.Millisecond * time.Duration(idleMs)
		if idleTime < idleTimeThreshold {
			continue // Do not reap active consumers
		}

		msg = info["name"]
		consumer, err := msg.ToString()
		if err != nil {
			return nil, apperr.NewError(apperr.CodeInternalServerError).
				WithCause(err).
				WithSummary("Failed to parse consumer name")
		}

		names = append(names, consumer)
	}

	return names, nil
}
