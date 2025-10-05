package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Client struct {
	db *gorm.DB
}

func NewClient(cfg ClientConfig) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to connect to database").
			WithCause(err)
	}

	return &Client{db: db}, nil
}
