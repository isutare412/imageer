package postgres

import "gorm.io/gorm"

func applyPagination[T any](q gorm.ChainInterface[T], limit, offset int) gorm.ChainInterface[T] {
	q = q.Offset(offset)
	q = q.Limit(limit)
	return q
}
