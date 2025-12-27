package postgres

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
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
	u, err := gorm.G[entity.User](r.db).
		Where(gen.User.ID.Eq(id)).
		First(ctx)
	if err != nil {
		return user, dbhelpers.WrapError(err, "Failed to find user %s", id)
	}

	return u.ToDomain(), nil
}

func (r *UserRepository) Upsert(ctx context.Context, user domain.User) (userCreated domain.User, err error) {
	u := entity.NewUser(user)

	if err := gorm.G[entity.User](r.db,
		clause.OnConflict{
			Columns: []clause.Column{{Name: gen.User.ID.Column().Name}},
			DoUpdates: clause.AssignmentColumns([]string{
				gen.User.UpdatedAt.Column().Name,
				gen.User.Nickname.Column().Name,
				gen.User.Email.Column().Name,
				gen.User.PhotoURL.Column().Name,
			}),
		}).
		Create(ctx, u); err != nil {
		return userCreated, dbhelpers.WrapError(err, "Failed to upsert user")
	}

	return u.ToDomain(), nil
}
