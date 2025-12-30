package postgres

import (
	"context"
	"fmt"

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

func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	user, err := gorm.G[entity.User](r.db).
		Where(gen.User.ID.Eq(id)).
		First(ctx)
	if err != nil {
		return domain.User{}, dbhelpers.WrapError(err, "Failed to find user %s", id)
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) Upsert(ctx context.Context, user domain.User) (domain.User, error) {
	usr := entity.NewUser(user)

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[entity.User](tx,
			clause.OnConflict{
				Columns: []clause.Column{{Name: gen.User.Email.Column().Name}},
				DoUpdates: clause.AssignmentColumns([]string{
					gen.User.UpdatedAt.Column().Name,
					gen.User.Nickname.Column().Name,
					gen.User.Email.Column().Name,
					gen.User.PhotoURL.Column().Name,
				}),
			}).
			Create(ctx, &usr); err != nil {
			return dbhelpers.WrapError(err, "Failed to upsert user")
		}

		usrFetched, err := gorm.G[entity.User](tx).
			Where(gen.User.Email.Eq(usr.Email)).
			First(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to fetch upserted user")
		}

		usr = usrFetched
		return nil
	}); err != nil {
		return domain.User{}, fmt.Errorf("during transaction: %w", err)
	}

	return usr.ToDomain(), nil
}
