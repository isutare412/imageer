package postgres_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/users"
)

func TestUserRepository_Upsert(t *testing.T) {
	postgresClient, mock := postgres.NewClientWithMock(t)
	userRepo := postgres.NewUserRepository(postgresClient)

	mock.ExpectBegin()
	mock.ExpectExec(
		`INSERT INTO "users" ("id","created_at","updated_at","role","nickname","email","photo_url") `+
			`VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT ("email") `+
			`DO UPDATE SET `+
			`"updated_at"="excluded"."updated_at",`+
			`"nickname"="excluded"."nickname",`+
			`"email"="excluded"."email",`+
			`"photo_url"="excluded"."photo_url"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "id" = $1 ORDER BY "users"."id" LIMIT $2`).
		WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.User]()).
			AddRow("user-1", time.Now(), time.Now(), users.RoleGuest, "nickname-1", "email-1", "photo-url-1"))
	mock.ExpectCommit()

	_, err := userRepo.Upsert(t.Context(), domain.User{})
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
