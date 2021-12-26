package user

import (
	"context"

	"gorm.io/gorm"
)

type Repo interface {
	Session(ctx context.Context) *gorm.DB
}
