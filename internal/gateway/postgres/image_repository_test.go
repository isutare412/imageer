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

func TestImageRepository_FindByID(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		imageRepo     *postgres.ImageRepository
		mock          sqlmock.Sqlmock

		id      string
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			id:   "image-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.imageRepo = postgres.NewImageRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectQuery(
					`SELECT * FROM "images" WHERE "id" = $1 ORDER BY "images"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Image]()).
						AddRow("image-1", time.Now(), time.Now(), "file-1.jpg", images.FormatJPEG,
							images.StateReady, "s3-key-1", "url-1", "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "projects"."id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "image_variants" WHERE "image_variants"."image_id" = $1`).
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
				_, err := tt.imageRepo.FindByID(ctx, tt.id)
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

func TestImageRepository_List(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		imageRepo     *postgres.ImageRepository
		mock          sqlmock.Sqlmock

		req     domain.ListImagesParams
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	updatedAtBefore := time.Now().Add(-10 * time.Minute)

	tests := []testSet{
		{
			name: "normal case",
			req: domain.ListImagesParams{
				Offset: lo.ToPtr(10),
				Limit:  lo.ToPtr(20),
				SearchFilter: domain.ImageSearchFilter{
					State:           lo.ToPtr(images.StateUploadPending),
					UpdatedAtBefore: &updatedAtBefore,
				},
				SortFilter: domain.ImageSortFilter{
					UpdatedAt: true,
					Direction: dbhelpers.SortDirectionDesc,
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.imageRepo = postgres.NewImageRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectQuery(
					`SELECT * FROM "images" WHERE "state" = $1 AND "updated_at" < $2 `+
						`ORDER BY "updated_at" DESC `+
						`LIMIT $3 OFFSET $4`).
					WithArgs(images.StateUploadPending, updatedAtBefore, 20, 10).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Image]()).
						AddRow("image-1", time.Now(), time.Now(), "file-1.jpg", images.FormatJPEG,
							images.StateUploadPending, "s3-key-1", "url-1", "project-1"))
				mock.ExpectQuery(`SELECT * FROM "projects" WHERE "projects"."id" = $1`).
					WithArgs("project-1").
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(`SELECT * FROM "image_variants" WHERE "image_variants"."image_id" = $1`).
					WithArgs("image-1").
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ImageVariant]()).
						AddRow("variant-1", time.Now(), time.Now(), images.FormatWebp,
							images.VariantStateReady, "s3-key-1", "url-1", "image-1", "preset-1"))
				mock.ExpectQuery(`SELECT * FROM "presets" WHERE "presets"."id" = $1`).
					WithArgs("preset-1").
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-name-1", false,
							images.FormatWebp, images.Quality(90), nil, nil, nil, nil, "project-1"))
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM "images" WHERE "state" = $1 AND "updated_at" < $2`).
					WithArgs(images.StateUploadPending, updatedAtBefore).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(5))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				_, err := tt.imageRepo.List(ctx, tt.req)
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

func TestImageRepository_Create(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		imageRepo     *postgres.ImageRepository
		mock          sqlmock.Sqlmock

		req     domain.Image
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.Image{
				FileName: "image-file-1.jpg",
				Format:   "jpeg",
				State:    "uploaded",
				S3Key:    "s3-key-1",
				Project: domain.ProjectReference{
					ID: "project-1",
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.imageRepo = postgres.NewImageRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "id" = $1 ORDER BY "projects"."id" LIMIT $2`).
					WithArgs(tt.req.Project.ID, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectExec(
					`INSERT INTO "images" ` +
						`("id","created_at","updated_at","file_name","format","state","s3_key","url","project_id") VALUES ` +
						`($1,$2,$3,$4,$5,$6,$7,$8,$9)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					`SELECT * FROM "images" WHERE "id" = $1 ORDER BY "images"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Image]()).
						AddRow("image-1", time.Now(), time.Now(), "file-1.jpg", images.FormatJPEG,
							images.StateReady, "s3-key-1", "url-1", "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "projects"."id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "image_variants" WHERE "image_variants"."image_id" = $1`).
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
				_, err := tt.imageRepo.Create(ctx, tt.req)
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

func TestImageRepository_Update(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		imageRepo     *postgres.ImageRepository
		mock          sqlmock.Sqlmock

		req     domain.UpdateImageRequest
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.UpdateImageRequest{
				ID:    "image-1",
				State: lo.ToPtr(images.StateReady),
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.imageRepo = postgres.NewImageRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(
					`UPDATE "images" SET "state"=$1,"updated_at"=NOW() WHERE "id" = $2`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					`SELECT * FROM "images" WHERE "id" = $1 ORDER BY "images"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Image]()).
						AddRow("image-1", time.Now(), time.Now(), "file-1.jpg", images.FormatJPEG,
							images.StateReady, "s3-key-1", "url-1", "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "projects"."id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "image_variants" WHERE "image_variants"."image_id" = $1`).
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
				_, err := tt.imageRepo.Update(ctx, tt.req)
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
