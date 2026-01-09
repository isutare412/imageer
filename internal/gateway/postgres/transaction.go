package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type txContextKey struct{}

type Transactioner struct {
	db *gorm.DB
}

func NewTransactioner(client *Client) *Transactioner {
	return &Transactioner{db: client.db}
}

func (t *Transactioner) BeginTx(ctx context.Context, opts ...*sql.TxOptions,
) (ctxWithTx context.Context, commit, rollback func() error) {
	tx := t.db.Begin(opts...)

	if txAny := ctx.Value(txContextKey{}); txAny != nil {
		panic("transaction is already in progress in the context")
	}

	ctxWithTx = context.WithValue(ctx, txContextKey{}, tx)

	commit = func() error {
		return tx.Commit().Error
	}
	rollback = func() error {
		return tx.Rollback().Error
	}

	return ctxWithTx, commit, rollback
}

func (t *Transactioner) WithTx(ctx context.Context, fn func(ctxWithTx context.Context) error) error {
	ctxWithTx, commit, rollback := t.BeginTx(ctx)

	defer func() {
		if v := recover(); v != nil {
			slog.Error("Panicked during transaction", "recover", v)

			if err := rollback(); err != nil {
				slog.Error("Failed to rollback transaction after panic", "error", err)
			}

			// Re-panic to propagate the panic up the call stack
			panic(v)
		}
	}()

	if ferr := fn(ctxWithTx); ferr != nil {
		ferr = fmt.Errorf("executing transactional function: %w", ferr)

		if rerr := rollback(); rerr != nil {
			rerr = fmt.Errorf("rolling back transaction: %w", rerr)
			ferr = errors.Join(ferr, rerr)
		}
		return ferr
	}

	if err := commit(); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}
	return nil
}

func GetTxOrDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok {
		return tx
	}
	return db
}
