package repository

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/repository/entity"
	"github.com/isutare412/imageer/pkg/apperr"
)

type Client struct {
	db *gorm.DB
}

func NewClient(cfg ClientConfig) (*Client, error) {
	db, err := gorm.Open(newDialector(cfg), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to connect to database").
			WithCause(err)
	}

	return &Client{db: db}, nil
}

func (c *Client) MigrateSchemas(ctx context.Context) error {
	if err := c.db.AutoMigrate(
		&entity.User{},
		&entity.Project{},
		&entity.Transformation{},
		&entity.ServiceAccount{},
		&entity.Image{},
	); err != nil {
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to migrate database schemas").
			WithCause(err)
	}
	return nil
}

func newDialector(cfg ClientConfig) gorm.Dialector {
	switch {
	case cfg.UseInMemory:
		return sqlite.Open("file::memory:?_foreign_key=true&mode=memory")

	default:
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
		return postgres.Open(dsn)
	}
}
