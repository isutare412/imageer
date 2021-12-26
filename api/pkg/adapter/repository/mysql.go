package repository

import (
	"fmt"

	"github.com/isutare412/imageer/api/pkg/config"
	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL(cfg *config.MySQLConfig) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Address)
	db, err := gorm.Open(mysqldriver.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("on open gorm DB: %w", err)
	}

	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", cfg.Database))
	db.Exec(fmt.Sprintf("USE %s;", cfg.Database))

	return &MySQL{
		db: db,
	}, nil
}
