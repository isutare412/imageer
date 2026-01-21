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

func TestProjectRepository_FindByID(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		projectRepo   *postgres.ProjectRepository
		mock          sqlmock.Sqlmock

		id      string
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			id:   "project-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.projectRepo = postgres.NewProjectRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT * FROM "projects" WHERE "id" = $1 ORDER BY "projects"."id" LIMIT $2`).
					WithArgs("project-1", 1).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(`SELECT * FROM "presets" WHERE "presets"."project_id" = $1`).
					WithArgs("project-1").
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-1", false, images.FormatWebp,
							images.Quality(90), images.FitCover, nil, 100, 100, "project-1"))
				mock.ExpectQuery(`SELECT COUNT(1) FROM "images" WHERE "project_id" = $1`).
					WithArgs("project-1").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(5))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				_, err := tt.projectRepo.FindByID(ctx, tt.id)
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

func TestProjectRepository_List(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		projectRepo   *postgres.ProjectRepository
		mock          sqlmock.Sqlmock

		req     domain.ListProjectsParams
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.ListProjectsParams{
				Offset: lo.ToPtr(20),
				Limit:  lo.ToPtr(20),
				SearchFilter: domain.ProjectSearchFilter{
					Name: lo.ToPtr("project-1"),
				},
				SortFilter: domain.ProjectSortFilter{
					UpdatedAt: true,
					Direction: dbhelpers.SortDirectionDesc,
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.projectRepo = postgres.NewProjectRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "name" = $1 `+
						`ORDER BY "updated_at" DESC `+
						`LIMIT $2 OFFSET $3`).
					WithArgs("project-1", 20, 20).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(`SELECT * FROM "presets" WHERE "presets"."project_id" = $1`).
					WithArgs("project-1").
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-1", false, images.FormatWebp,
							images.Quality(90), images.FitCover, nil, 100, 100, "project-1"))
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM "projects" WHERE "name" = $1`).
					WithArgs(tt.req.SearchFilter.Name).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(5))
				mock.ExpectQuery(
					`SELECT "project_id",COUNT(1) AS count FROM "images" WHERE "project_id" = $1 GROUP BY "project_id"`).
					WithArgs("project-1").
					WillReturnRows(sqlmock.NewRows([]string{"project_id", "count"}).AddRow("project-1", 3))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				_, err := tt.projectRepo.List(ctx, tt.req)
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

func TestProjectRepository_Create(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		projectRepo   *postgres.ProjectRepository
		mock          sqlmock.Sqlmock

		req     domain.Project
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.Project{
				Name: "project-1",
				Presets: []domain.Preset{
					{
						Name:    "preset-name-1",
						Default: false,
						Quality: images.Quality(90),
						Format:  images.FormatWebp,
						Fit:     lo.ToPtr(images.FitCover),
						Anchor:  lo.ToPtr(images.AnchorSmart),
						Width:   lo.ToPtr[int64](100),
						Height:  lo.ToPtr[int64](100),
					},
					{
						Name:    "preset-name-2",
						Default: false,
						Format:  images.FormatWebp,
						Quality: images.Quality(90),
						Fit:     lo.ToPtr(images.FitCover),
						Anchor:  lo.ToPtr(images.AnchorSmart),
						Width:   lo.ToPtr[int64](200),
						Height:  lo.ToPtr[int64](200),
					},
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.projectRepo = postgres.NewProjectRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(
					`INSERT INTO "projects" ("id","created_at","updated_at","name") VALUES ($1,$2,$3,$4)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`INSERT INTO "presets" ` +
						`("id","created_at","updated_at","name","default","format","quality","fit","anchor","width","height","project_id") VALUES ` +
						`($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12),($13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24) ` +
						`ON CONFLICT ("id") DO UPDATE SET "project_id"="excluded"."project_id"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				_, err := tt.projectRepo.Create(ctx, tt.req)
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

func TestProjectRepository_Update(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		projectRepo   *postgres.ProjectRepository
		mock          sqlmock.Sqlmock

		req     domain.UpdateProjectRequest
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.UpdateProjectRequest{
				ID:   "project-1",
				Name: lo.ToPtr("project-1"),
				Presets: []domain.UpsertPresetRequest{
					{
						// Update request
						ID:     lo.ToPtr("preset-1"),
						Name:   lo.ToPtr("preset-1"),
						Width:  lo.ToPtr[int64](100),
						Height: lo.ToPtr[int64](100),
					},
					{
						// Create request
						Name:    lo.ToPtr("preset-2"),
						Default: lo.ToPtr(true),
						Width:   lo.ToPtr[int64](100),
						Height:  lo.ToPtr[int64](100),
					},
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.projectRepo = postgres.NewProjectRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(
					`UPDATE "presets" SET "name"=$1,"width"=$2,"height"=$3,"updated_at"=NOW() WHERE "id" = $4`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`INSERT INTO "presets" ` +
						`("id","created_at","updated_at","name","default","format","quality","fit","anchor","width","height","project_id") VALUES ` +
						`($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`DELETE FROM "presets" WHERE "project_id" = $1 AND "id" NOT IN ($2,$3)`).
					WithArgs("project-1", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`UPDATE "projects" SET "name"=$1,"updated_at"=NOW() WHERE "id" = $2`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "id" = $1 ORDER BY "projects"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "presets" WHERE "presets"."project_id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-1", false, images.FormatWebp,
							images.Quality(90), images.FitCover, nil, 100, 100, "project-1").
						AddRow("preset-2", time.Now(), time.Now(), "preset-2", false, images.FormatWebp,
							images.Quality(90), images.FitCover, nil, 100, 100, "project-1"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				_, err := tt.projectRepo.Update(ctx, tt.req)
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

func TestProjectRepository_Delete(t *testing.T) {
	type testSet struct {
		name          string // description of this test case
		transactioner *postgres.Transactioner
		projectRepo   *postgres.ProjectRepository
		mock          sqlmock.Sqlmock

		id      string
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			id:   "project-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.projectRepo = postgres.NewProjectRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "projects" WHERE "id" = $1`).
					WithArgs(tt.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				err := tt.projectRepo.Delete(ctx, tt.id)
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
