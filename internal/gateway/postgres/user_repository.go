package postgres

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/trace"
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
	ctx, span := trace.StartSpan(ctx, "postgres.UserRepository.FindByID")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	user, err := gorm.G[entity.User](tx).
		Where(gen.User.ID.Eq(id)).
		First(ctx)
	if err != nil {
		return domain.User{}, dbhelpers.WrapGORMError(err, "Failed to find user %s", id)
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) Upsert(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, span := trace.StartSpan(ctx, "postgres.UserRepository.Upsert")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	usr := entity.NewUser(user)
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
		return domain.User{}, dbhelpers.WrapGORMError(err, "Failed to upsert user")
	}

	usr, err := gorm.G[entity.User](tx).
		Where(gen.User.Email.Eq(usr.Email)).
		First(ctx)
	if err != nil {
		return domain.User{}, dbhelpers.WrapGORMError(err, "Failed to fetch upserted user")
	}

	return usr.ToDomain(), nil
}
