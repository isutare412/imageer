package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	sloggorm "github.com/orandin/slog-gorm"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClientWithMock(t *testing.T) (*Client, sqlmock.Sqlmock) {
	db, mock, err := dbhelpers.NewSQLMock()
	require.NoError(t, err)

	slogAdapter := sloggorm.New(
		sloggorm.WithTraceAll(),
	)

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		TranslateError: true,
		Logger:         slogAdapter,
	})
	require.NoError(t, err)

	return &Client{db: gdb}, mock
}
