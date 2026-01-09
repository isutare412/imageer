package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sloggorm "github.com/orandin/slog-gorm"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/pkg/dbhelpers"
)

func NewClientWithMock(t *testing.T) (*Client, *Transactioner, sqlmock.Sqlmock) {
	db, mock, err := dbhelpers.NewSQLMock()
	require.NoError(t, err)

	slogAdapter := sloggorm.New(
	// sloggorm.WithTraceAll(),
	)

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		TranslateError: true,
		Logger:         slogAdapter,
	})
	require.NoError(t, err)

	client := &Client{db: gdb}
	return client, NewTransactioner(client), mock
}
