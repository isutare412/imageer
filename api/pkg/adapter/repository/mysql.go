package repository

import (
	"context"
	"fmt"

	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/api/pkg/config"
	"github.com/isutare412/imageer/api/pkg/core/user"
)

type MySQL struct {
	db *gorm.DB
}

func (r *MySQL) Session(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func NewMySQL(cfg *config.MySQLConfig) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Address, cfg.Database)
	db, err := gorm.Open(mysqldriver.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("on open gorm DB: %w", err)
	}

	if err := db.AutoMigrate(
		&user.User{},
	); err != nil {
		return nil, fmt.Errorf("on migrate entities: %w", err)
	}

	return &MySQL{
		db: db,
	}, nil
}
