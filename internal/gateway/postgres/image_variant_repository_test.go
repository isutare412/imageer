package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/images"
)

func TestImageVariantRepository_Create(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		imageVarRepo  *postgres.ImageVariantRepository
		mock          sqlmock.Sqlmock

		req     domain.ImageVariant
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.ImageVariant{
				Format:  images.FormatWebp,
				State:   images.VariantStateReady,
				S3Key:   "s3-key-1",
				ImageID: "image-1",
				Preset: domain.PresetReference{
					ID: "preset-1",
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.imageVarRepo = postgres.NewImageVariantRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectQuery(
					`SELECT * FROM "images" WHERE "id" = $1 ORDER BY "images"."id" LIMIT $2`).
					WithArgs(tt.req.ImageID, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Image]()).
						AddRow("image-1", time.Now(), time.Now(), "file-1.jpg", images.FormatJPEG,
							images.StateReady, "s3-key-1", "url-1", "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "presets" WHERE "id" = $1 ORDER BY "presets"."id" LIMIT $2`).
					WithArgs(tt.req.Preset.ID, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-name-1", false,
							images.FormatWebp, images.Quality(90), nil, nil, nil, nil, "project-1"))
				mock.ExpectExec(
					`INSERT INTO "image_variants" ` +
						`("id","created_at","updated_at","format","state","s3_key","url","image_id","preset_id") VALUES ` +
						`($1,$2,$3,$4,$5,$6,$7,$8,$9)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					`SELECT * FROM "image_variants" WHERE "id" = $1 ORDER BY "image_variants"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ImageVariant]()).
						AddRow("variant-1", time.Now(), time.Now(), images.FormatWebp,
							images.VariantStateReady, "s3-key-1", "url-1", "image-1", "preset-1"))
				mock.ExpectQuery(
					`SELECT * FROM "presets" WHERE "presets"."id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-name-1", false,
							images.FormatWebp, images.Quality(90), nil, nil, nil, nil, "project-1"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				_, err := tt.imageVarRepo.Create(ctx, tt.req)
				return err
			})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			err = tt.mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestImageVariantRepository_Update(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		imageVarRepo  *postgres.ImageVariantRepository
		mock          sqlmock.Sqlmock

		req     domain.UpdateImageVariantRequest
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.UpdateImageVariantRequest{
				ID:    "variant-1",
				State: lo.ToPtr(images.VariantStateReady),
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.imageVarRepo = postgres.NewImageVariantRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(
					`UPDATE "image_variants" SET "state"=$1,"updated_at"=NOW() WHERE "id" = $2`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					`SELECT * FROM "image_variants" WHERE "id" = $1 ORDER BY "image_variants"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ImageVariant]()).
						AddRow("variant-1", time.Now(), time.Now(), images.FormatWebp,
							images.VariantStateReady, "s3-key-1", "url-1", "image-1", "preset-1"))
				mock.ExpectQuery(
					`SELECT * FROM "presets" WHERE "presets"."id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-name-1", false,
							images.FormatWebp, images.Quality(90), nil, nil, nil, nil, "project-1"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				_, err := tt.imageVarRepo.Update(ctx, tt.req)
				return err
			})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			err = tt.mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
