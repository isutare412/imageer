package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"

	"github.com/go-sql-driver/mysql"
	"github.com/isutare412/imageer/api-server/pkg/config"
	"github.com/isutare412/imageer/api-server/pkg/core/user"
)

type MySQL struct {
	db *gorm.DB
}

func (r *MySQL) Session(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *MySQL) IsErrDuplicate(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func (r *MySQL) IsErrNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func NewMySQL(cfg *config.MySQLConfig) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Address, cfg.Database)

	useColor := false
	if config.IsDevelopmentMode() {
		useColor = true
	}
	logger := glog.New(
		&gormLogger{},
		glog.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  glog.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  useColor,
		},
	)

	db, err := gorm.Open(mysqldriver.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&user.User{},
	); err != nil {
		return nil, err
	}

	return &MySQL{
		db: db,
	}, nil
}
