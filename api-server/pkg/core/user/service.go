package user

import (
	"context"
	"fmt"

	"github.com/isutare412/imageer/api-server/pkg/core/auth"
)

type Service interface {
	Create(ctx context.Context, user *User, password string) (int64, error)
	GetByEmailPassword(ctx context.Context, email, password string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	UpdateCredit(ctx context.Context, id int64, delta int64) (*User, error)
}

type service struct {
	repo    Repo
	authSvc auth.Service
}

func (s *service) Create(ctx context.Context, user *User, password string) (id int64, err error) {
	user.Privilege = PrivilegeUser

	hashed, err := s.authSvc.Hash(password)
	if err != nil {
		return 0, fmt.Errorf("on create user: %w", err)
	}
	user.Password = hashed

	tx := s.repo.Session(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	res := tx.Where("email = ?", user.Email).Find(&User{})
	if err := res.Error; err != nil {
		return 0, fmt.Errorf("on create user: %w", err)
	} else if res.RowsAffected > 0 {
		return 0, fmt.Errorf("email[%s] duplicated: %w", user.Email, ErrDuplicate)
	}

	err = tx.Create(user).Error
	if err != nil {
		return 0, fmt.Errorf("on create user: %w", err)
	}
	return user.ID, nil
}

func (s *service) GetByEmailPassword(ctx context.Context, email, password string) (*User, error) {
	db := s.repo.Session(ctx)

	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if s.repo.IsErrNotFound(err) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("on get user by email: %w", err)
	}

	if ok := s.authSvc.Compare(password, user.Password); !ok {
		return nil, ErrPasswordNotCorrect
	}
	return &user, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	db := s.repo.Session(ctx)

	var user User
	err := db.First(&user, id).Error
	if s.repo.IsErrNotFound(err) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("on get by id: %w", err)
	}
	return &user, nil
}

func (s *service) UpdateCredit(ctx context.Context, id int64, delta int64) (user *User, err error) {
	tx := s.repo.Session(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.First(user, id).Error
	if s.repo.IsErrNotFound(err) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("on update credit: %w", err)
	}

	newCredit := user.Credit + delta
	if newCredit < 0 {
		newCredit = 0
	}
	err = tx.Model(&user).Update("credit", newCredit).Error
	if err != nil {
		return nil, fmt.Errorf("on update credit: %w", err)
	}
	user.Credit = newCredit

	return user, nil
}

func NewService(repo Repo, authSvc auth.Service) Service {
	return &service{
		repo:    repo,
		authSvc: authSvc,
	}
}
