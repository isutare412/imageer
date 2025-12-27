package dbhelpers

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/isutare412/imageer/pkg/apperr"
)

const (
	summaryResourceNotFound    = "Resource not found"
	summaryResourceConflict    = "Resource conflict"
	summaryInternalServerError = "Internal server error"
)

func WrapError(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return apperr.NewError(apperr.CodeNotFound).
			WithSummary(summaryResourceNotFound).
			WithDetail(msg, args...).
			WithCause(err)
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return apperr.NewError(apperr.CodeConflict).
			WithSummary(summaryResourceConflict).
			WithDetail(msg, args...).
			WithCause(err)
	default:
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary(summaryInternalServerError).
			WithDetail(msg, args...).
			WithCause(err)
	}
}

func ColumnNamesFor[T any]() []string {
	itf := reflect.New(reflect.TypeFor[T]()).Interface()
	s, err := schema.Parse(itf, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		panic(fmt.Errorf("unexpected schema parsing error: %v", err))
	}

	columns := lo.FilterMap(s.Fields, func(f *schema.Field, _ int) (string, bool) {
		return f.DBName, f.DBName != ""
	})
	return columns
}
