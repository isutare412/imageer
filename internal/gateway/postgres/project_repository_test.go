package postgres_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

func TestProjectRepository_FindByID(t *testing.T) {
	type testSet struct {
		name        string // description of this test case
		projectRepo *postgres.ProjectRepository
		mock        sqlmock.Sqlmock

		id      string
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			id:   "project-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.projectRepo = postgres.NewProjectRepository(postgresClient)
				tt.mock = mock

				mock.ExpectQuery(`SELECT * FROM "projects" WHERE "id" = $1 ORDER BY "projects"."id" LIMIT $2`).
					WithArgs("project-1", 1).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(`SELECT * FROM "transformations" WHERE "transformations"."project_id" = $1`).
					WithArgs("project-1").
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Transformation]()).
						AddRow("trans-1", time.Now(), time.Now(), "trans-1", false, 100, 100, "project-1"))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.projectRepo.FindByID(t.Context(), tt.id)
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
		name        string // description of this test case
		projectRepo *postgres.ProjectRepository
		mock        sqlmock.Sqlmock

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
				postgresClient, mock := postgres.NewClientWithMock(t)
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
				mock.ExpectQuery(`SELECT * FROM "transformations" WHERE "transformations"."project_id" = $1`).
					WithArgs("project-1").
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Transformation]()).
						AddRow("trans-1", time.Now(), time.Now(), "trans-1", false, 100, 100, "project-1"))
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM "projects" WHERE "name" = $1`).
					WithArgs(tt.req.SearchFilter.Name).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(5))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.projectRepo.List(t.Context(), tt.req)
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
		name        string // description of this test case
		projectRepo *postgres.ProjectRepository
		mock        sqlmock.Sqlmock

		req     domain.Project
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.Project{
				Name: "project-1",
				Transformations: []domain.Transformation{
					{
						Name:    "trans-name-1",
						Default: false,
						Width:   100,
						Height:  100,
					},
					{
						Name:    "trans-name-2",
						Default: false,
						Width:   200,
						Height:  200,
					},
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.projectRepo = postgres.NewProjectRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(
					`INSERT INTO "projects" ("id","created_at","updated_at","name") VALUES ($1,$2,$3,$4)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`INSERT INTO "transformations" ` +
						`("id","created_at","updated_at","name","default","width","height","project_id") VALUES ` +
						`($1,$2,$3,$4,$5,$6,$7,$8),($9,$10,$11,$12,$13,$14,$15,$16) ` +
						`ON CONFLICT ("id") ` +
						`DO UPDATE SET "project_id"="excluded"."project_id"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.projectRepo.Create(t.Context(), tt.req)
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
		name        string // description of this test case
		projectRepo *postgres.ProjectRepository
		mock        sqlmock.Sqlmock

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
				Transformations: []domain.UpsertTransformationRequest{
					{
						// Update request
						ID:     lo.ToPtr("trans-1"),
						Name:   lo.ToPtr("trans-1"),
						Width:  lo.ToPtr[int64](100),
						Height: lo.ToPtr[int64](100),
					},
					{
						// Create request
						Name:    lo.ToPtr("trans-2"),
						Default: lo.ToPtr(true),
						Width:   lo.ToPtr[int64](100),
						Height:  lo.ToPtr[int64](100),
					},
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.projectRepo = postgres.NewProjectRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(
					`UPDATE "transformations" SET "name"=$1,"width"=$2,"height"=$3,"updated_at"=$4 WHERE "id" = $5`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`INSERT INTO "transformations" ` +
						`("id","created_at","updated_at","name","default","width","height","project_id") VALUES ` +
						`($1,$2,$3,$4,$5,$6,$7,$8)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`DELETE FROM "transformations" WHERE "id" NOT IN ($1,$2)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`UPDATE "projects" SET "name"=$1,"updated_at"=$2 WHERE "id" = $3`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "id" = $1 ORDER BY "projects"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-1"))
				mock.ExpectQuery(
					`SELECT * FROM "transformations" WHERE "transformations"."project_id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Transformation]()).
						AddRow("trans-1", time.Now(), time.Now(), "trans-1", false,
							int64(100), int64(100), "project-1").
						AddRow("trans-2", time.Now(), time.Now(), "trans-2", true,
							int64(100), int64(100), "project-1"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.projectRepo.Update(t.Context(), tt.req)
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
		name        string // description of this test case
		projectRepo *postgres.ProjectRepository
		mock        sqlmock.Sqlmock

		id      string
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			id:   "project-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
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

			err := tt.projectRepo.Delete(t.Context(), tt.id)
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
