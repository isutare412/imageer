package postgres_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres"
)

func TestUserRepository_Upsert(t *testing.T) {
	postgresClient, mock := postgres.NewClientWithMock(t)
	userRepo := postgres.NewUserRepository(postgresClient)

	mock.ExpectBegin()
	mock.ExpectExec(
		`INSERT INTO "users" ("id","created_at","updated_at","role","nickname","email","photo_url") `+
			`VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("id") `+
			`DO UPDATE SET `+
			`"updated_at"="excluded"."updated_at",`+
			`"nickname"="excluded"."nickname",`+
			`"email"="excluded"."email",`+
			`"photo_url"="excluded"."photo_url"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err := userRepo.Upsert(t.Context(), domain.User{})
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
