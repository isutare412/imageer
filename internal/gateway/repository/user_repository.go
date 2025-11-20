package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/repository/entity"
	"github.com/isutare412/imageer/pkg/apperr"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(client *Client) *UserRepository {
	return &UserRepository{
		db: client.db,
	}
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (user domain.User, err error) {
	u, err := gorm.G[entity.User](r.db).First(ctx)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return user, apperr.NewError(apperr.CodeNotFound).
			WithSummary("user not found by id %s", id).
			WithCause(err)
	case err != nil:
		return user, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to find user by id %s", id).
			WithCause(err)
	}

	return u.ToDomain(), nil
}

func (r *UserRepository) Upsert(ctx context.Context, user domain.User) (userCreated domain.User, err error) {
	u := entity.NewUser(user)

	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"updated_at",
				"nickname",
				"email",
				"photo_url",
			}),
		}).Create(&u).Error; err != nil {
		return userCreated, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to upsert user").
			WithCause(err)
	}

	return u.ToDomain(), nil
}
