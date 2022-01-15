package user

import (
	"context"
	"fmt"

	"gorm.io/gorm"

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

func (s *service) Create(ctx context.Context, user *User, password string) (int64, error) {
	hashed, err := s.authSvc.Hash(password)
	if err != nil {
		return 0, fmt.Errorf("on create user: %w", err)
	}
	user.Password = hashed

	var newID int64
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("on create user: %w", err)
		}

		newID = user.ID
		return nil
	}); s.repo.IsErrDuplicate(err) {
		return 0, ErrDuplicate
	} else if err != nil {
		return 0, fmt.Errorf("on create user: %w", err)
	}

	return newID, nil
}

func (s *service) GetByEmailPassword(ctx context.Context, email, password string) (*User, error) {
	var user User
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("email = ?", email).First(&user).Error; err != nil {
			return fmt.Errorf("on get user by email: %w", err)
		}
		return nil
	}); s.repo.IsErrNotFound(err) {
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
	var user User
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return fmt.Errorf("on get by id: %w", err)
		}
		return nil
	}); s.repo.IsErrNotFound(err) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("on get by id: %w", err)
	}

	return &user, nil
}

func (s *service) UpdateCredit(ctx context.Context, id int64, delta int64) (*User, error) {
	var user User
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return fmt.Errorf("on update credit: %w", err)
		}

		newCredit := user.Credit + delta
		if newCredit < 0 {
			newCredit = 0
		}
		if err := tx.Model(&user).Update("credit", newCredit).Error; err != nil {
			return fmt.Errorf("on update credit: %w", err)
		}
		user.Credit = newCredit
		return nil
	}); err != nil {
		return nil, fmt.Errorf("on update credit: %w", err)
	}
	return &user, nil
}

func NewService(repo Repo, authSvc auth.Service) Service {
	return &service{
		repo:    repo,
		authSvc: authSvc,
	}
}
