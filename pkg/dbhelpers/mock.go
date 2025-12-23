package dbhelpers

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewSQLMock() (*sql.DB, sqlmock.Sqlmock, error) {
	return sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
}
