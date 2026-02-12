package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres"
)

func TestImageProcessingLogRepository_Create(t *testing.T) {
	type testSet struct {
		name             string // description of this test case
		transactioner    *postgres.Transactioner
		imageProcLogRepo *postgres.ImageProcessingLogRepository
		mock             sqlmock.Sqlmock

		req     domain.ImageProcessingLog
		setup   func(t *testing.T, tt *testSet)
		wantErr bool
	}

	tests := []testSet{
		{
			name: "normal case",
			req: domain.ImageProcessingLog{
				IsSuccess:      false,
				ErrorCode:      new(42),
				ErrorMessage:   new("Some error occurred"),
				ElapsedTime:    812 * time.Millisecond,
				ImageVariantID: "variant-1",
			},
			setup: func(t *testing.T, tt *testSet) {
				postgresClient, transactioner, mock := postgres.NewClientWithMock(t)
				tt.transactioner = transactioner
				tt.imageProcLogRepo = postgres.NewImageProcessingLogRepository(postgresClient)
				tt.mock = mock

				mock.ExpectBegin()
				mock.ExpectQuery(
					`INSERT INTO "image_processing_logs" ` +
						`("created_at","is_success","error_code","error_message","elapsed_time_millis","image_variant_id") VALUES ` +
						`($1,$2,$3,$4,$5,$6) RETURNING "id"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, &tt)

			err := tt.transactioner.WithTx(t.Context(), func(ctx context.Context) error {
				_, err := tt.imageProcLogRepo.Create(ctx, tt.req)
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
