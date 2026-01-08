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
	"github.com/isutare412/imageer/pkg/images"
)

func TestPresetRepository_FindByName(t *testing.T) {
	type testSet struct {
		name       string // description of this test case
		presetRepo *postgres.PresetRepository
		mock       sqlmock.Sqlmock

		projectID  string
		presetName string
		setup      func(t *testing.T, tt *testSet)
		wantErr    bool
	}

	tests := []testSet{
		{
			name:       "normal case",
			projectID:  "project-1",
			presetName: "preset-name-1",
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.presetRepo = postgres.NewPresetRepository(postgresClient)
				tt.mock = mock

				mock.ExpectQuery(
					`SELECT * FROM "presets" WHERE "project_id" = $1 AND "name" = $2 `+
						`ORDER BY "presets"."id" LIMIT $3`).
					WithArgs("project-1", "preset-name-1", 1).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-name-1", false,
							images.FormatWebp, images.Quality(90), nil, nil, nil, nil, "project-1"))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.presetRepo.FindByName(t.Context(), tt.projectID, tt.presetName)
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

func TestPresetRepository_List(t *testing.T) {
	type testSet struct {
		name       string // description of this test case
		presetRepo *postgres.PresetRepository
		mock       sqlmock.Sqlmock

		req     domain.ListPresetsParams
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.ListPresetsParams{
				Offset: lo.ToPtr(20),
				Limit:  lo.ToPtr(20),
				SearchFilter: domain.PresetSearchFilter{
					ProjectID: lo.ToPtr("project-1"),
					Names:     []string{"preset-name-1", "preset-name-2"},
				},
				SortFilter: domain.PresetSortFilter{
					UpdatedAt: true,
					Direction: dbhelpers.SortDirectionDesc,
				},
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, mock := postgres.NewClientWithMock(t)
				tt.presetRepo = postgres.NewPresetRepository(postgresClient)
				tt.mock = mock

				mock.ExpectQuery(
					`SELECT * FROM "presets" WHERE "project_id" = $1 AND "name" IN ($2,$3) `+
						`ORDER BY "updated_at" DESC LIMIT $4 OFFSET $5`).
					WithArgs("project-1", "preset-name-1", "preset-name-2", 20, 20).
					WillReturnRows(sqlmock.NewRows(dbhelpers.ColumnNamesFor[entity.Preset]()).
						AddRow("preset-1", time.Now(), time.Now(), "preset-name-1", false,
							images.FormatWebp, images.Quality(90), nil, nil, nil, nil, "project-1").
						AddRow("preset-2", time.Now(), time.Now(), "preset-name-2", false,
							images.FormatWebp, images.Quality(90), nil, nil, nil, nil, "project-1"))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			_, err := tt.presetRepo.List(t.Context(), tt.req)
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
