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
	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

func TestServiceAccountRepository_FindByID(t *testing.T) {
	type testSet struct {
		name           string // description of this test case
		svcAccountRepo *postgres.ServiceAccountRepository
		mock           sqlmock.Sqlmock

		id      string
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			id:   "account-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.svcAccountRepo = postgres.NewServiceAccountRepository(postgresClient)
				tt.mock = mock

				mock.ExpectQuery(
					`SELECT * FROM "service_accounts" WHERE "id" = $1 ORDER BY "service_accounts"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccount]()).
						AddRow("account-1", time.Now(), time.Now(), "account-name-1",
							serviceaccounts.AccessScopeFull, time.Now(), "test-api-key"))
				mock.ExpectQuery(
					`SELECT * FROM "service_account_projects" WHERE ` +
						`"service_account_projects"."service_account_id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccountProject]()).
						AddRow("account-1", "project-1").
						AddRow("account-1", "project-2"))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "projects"."id" IN ($1,$2)`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-name-1").
						AddRow("project-2", time.Now(), time.Now(), "project-name-2"))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.svcAccountRepo.FindByID(t.Context(), tt.id)
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

func TestServiceAccountRepository_FindByAPIKeyHash(t *testing.T) {
	type testSet struct {
		name           string // description of this test case
		svcAccountRepo *postgres.ServiceAccountRepository
		mock           sqlmock.Sqlmock

		hash    string
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			hash: "test-hash-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.svcAccountRepo = postgres.NewServiceAccountRepository(postgresClient)
				tt.mock = mock

				mock.ExpectQuery(
					`SELECT * FROM "service_accounts" WHERE "api_key_hash" = $1 ` +
						`ORDER BY "service_accounts"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccount]()).
						AddRow("account-1", time.Now(), time.Now(), "account-name-1",
							serviceaccounts.AccessScopeFull, time.Now(), "test-api-key"))
				mock.ExpectQuery(
					`SELECT * FROM "service_account_projects" WHERE ` +
						`"service_account_projects"."service_account_id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccountProject]()).
						AddRow("account-1", "project-1").
						AddRow("account-1", "project-2"))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "projects"."id" IN ($1,$2)`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-name-1").
						AddRow("project-2", time.Now(), time.Now(), "project-name-2"))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.svcAccountRepo.FindByAPIKeyHash(t.Context(), tt.hash)
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

func TestServiceAccountRepository_List(t *testing.T) {
	type testSet struct {
		name           string // description of this test case
		svcAccountRepo *postgres.ServiceAccountRepository
		mock           sqlmock.Sqlmock

		req     domain.ListServiceAccountsParams
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.ListServiceAccountsParams{
				Offset: lo.ToPtr(20),
				Limit:  lo.ToPtr(20),
				SearchFilter: domain.ServiceAccountSearchFilter{
					Name: lo.ToPtr("account-1"),
				},
				SortFilter: domain.ServiceAccountSortFilter{
					UpdatedAt: true,
					Direction: dbhelpers.SortDirectionDesc,
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.svcAccountRepo = postgres.NewServiceAccountRepository(postgresClient)
				tt.mock = mock

				mock.ExpectQuery(
					`SELECT * FROM "service_accounts" WHERE "name" = $1 `+
						`ORDER BY "updated_at" DESC LIMIT $2 OFFSET $3`).
					WithArgs("account-1", 20, 20).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccount]()).
						AddRow("account-1", time.Now(), time.Now(), "account-name-1",
							serviceaccounts.AccessScopeFull, time.Now(), "test-api-key"))
				mock.ExpectQuery(
					`SELECT * FROM "service_account_projects" WHERE ` +
						`"service_account_projects"."service_account_id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccountProject]()).
						AddRow("account-1", "project-1").
						AddRow("account-1", "project-2"))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "projects"."id" IN ($1,$2)`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-name-1").
						AddRow("project-2", time.Now(), time.Now(), "project-name-2"))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.svcAccountRepo.List(t.Context(), tt.req)
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

func TestServiceAccountRepository_Create(t *testing.T) {
	type testSet struct {
		name           string // description of this test case
		svcAccountRepo *postgres.ServiceAccountRepository
		mock           sqlmock.Sqlmock

		req     domain.ServiceAccount
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.ServiceAccount{
				ExpireAt:    lo.ToPtr(time.Now().Add(24 * time.Hour)),
				Name:        "account-name-1",
				AccessScope: serviceaccounts.AccessScopeFull,
				APIKeyHash:  "test-api-key",
				Projects: []domain.ProjectReference{
					{ID: "project-1"},
					{ID: "project-2"},
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.svcAccountRepo = postgres.NewServiceAccountRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(
					`INSERT INTO "service_accounts" ` +
						`("id","created_at","updated_at","name","access_scope","expire_at","api_key_hash") VALUES ` +
						`($1,$2,$3,$4,$5,$6,$7)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`INSERT INTO "service_account_projects" `+
					`("service_account_id","project_id") VALUES `+
					`($1,$2),($3,$4)`).
					WithArgs(sqlmock.AnyArg(), "project-1", sqlmock.AnyArg(), "project-2").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					`SELECT * FROM "service_accounts" WHERE "id" = $1 ORDER BY "service_accounts"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccount]()).
						AddRow("account-1", time.Now(), time.Now(),
							tt.req.Name, tt.req.AccessScope, tt.req.ExpireAt, tt.req.APIKeyHash))
				mock.ExpectQuery(
					`SELECT * FROM "service_account_projects" WHERE ` +
						`"service_account_projects"."service_account_id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccountProject]()).
						AddRow("account-1", "project-1").
						AddRow("account-1", "project-2"))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "projects"."id" IN ($1,$2)`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-name-1").
						AddRow("project-2", time.Now(), time.Now(), "project-name-2"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.svcAccountRepo.Create(t.Context(), tt.req)
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

func TestServiceAccountRepository_Update(t *testing.T) {
	type testSet struct {
		name           string // description of this test case
		svcAccountRepo *postgres.ServiceAccountRepository
		mock           sqlmock.Sqlmock

		req     domain.UpdateServiceAccountRequest
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.UpdateServiceAccountRequest{
				ID:          "account-1",
				Name:        lo.ToPtr("account-name-1"),
				AccessScope: lo.ToPtr(serviceaccounts.AccessScopeProject),
				ProjectIDs: []string{
					"project-1",
					"project-2",
				},
				ExpireAt: lo.ToPtr(time.Now().Add(time.Hour)),
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.svcAccountRepo = postgres.NewServiceAccountRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(
					`UPDATE "service_accounts" SET ` +
						`"name"=$1,"access_scope"=$2,"expire_at"=$3 WHERE "id" = $4`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`DELETE FROM "service_account_projects" WHERE "service_account_id" = $1`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(
					`INSERT INTO "service_account_projects" ` +
						`("service_account_id","project_id") VALUES ($1,$2),($3,$4)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					`SELECT * FROM "service_accounts" WHERE "id" = $1 ORDER BY "service_accounts"."id" LIMIT $2`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccount]()).
						AddRow("account-1", time.Now(), time.Now(),
							tt.req.Name, tt.req.AccessScope, tt.req.ExpireAt, "test-api-key"))
				mock.ExpectQuery(
					`SELECT * FROM "service_account_projects" WHERE ` +
						`"service_account_projects"."service_account_id" = $1`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.ServiceAccountProject]()).
						AddRow("account-1", "project-1").
						AddRow("account-1", "project-2"))
				mock.ExpectQuery(
					`SELECT * FROM "projects" WHERE "projects"."id" IN ($1,$2)`).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Project]()).
						AddRow("project-1", time.Now(), time.Now(), "project-name-1").
						AddRow("project-2", time.Now(), time.Now(), "project-name-2"))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.svcAccountRepo.Update(t.Context(), tt.req)
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

func TestServiceAccountRepository_Delete(t *testing.T) {
	type testSet struct {
		name           string // description of this test case
		svcAccountRepo *postgres.ServiceAccountRepository
		mock           sqlmock.Sqlmock

		id      string
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			id:   "account-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.svcAccountRepo = postgres.NewServiceAccountRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "service_accounts" WHERE "id" = $1`).
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

			err := tt.svcAccountRepo.Delete(t.Context(), tt.id)
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
