package user

import (
	"context"
	"fmt"

	"github.com/isutare412/imageer/api-server/pkg/core/encrypt"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, user *User, password string) (int64, error)
	GetByEmailPw(ctx context.Context, email, password string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	UpdateCredit(ctx context.Context, id int64, delta int64) (*User, error)
}

type service struct {
	repo   Repo
	ecrSvc encrypt.Service
}

func (s *service) Create(ctx context.Context, user *User, password string) (int64, error) {
	hashed, err := s.ecrSvc.Hash(password)
	if err != nil {
		return 0, fmt.Errorf("on hash password: %w", err)
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
		return 0, fmt.Errorf("on transaction: %w", err)
	}

	return newID, nil
}

func (s *service) GetByEmailPw(ctx context.Context, email, password string) (*User, error) {
	var user User
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("email = ?", email).Find(&user).Error; err != nil {
			return fmt.Errorf("on find user with email: %w", err)
		}
		return nil
	}); s.repo.IsErrNotFound(err) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("on transaction: %w", err)
	}

	if ok := s.ecrSvc.Compare(password, user.Password); !ok {
		return nil, ErrPasswordNotCorrect
	}

	return &user, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return fmt.Errorf("on first user: %w", err)
		}
		return nil
	}); s.repo.IsErrNotFound(err) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("on transaction: %w", err)
	}

	return &user, nil
}

func (s *service) UpdateCredit(ctx context.Context, id int64, delta int64) (*User, error) {
	var user User
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return fmt.Errorf("on first user: %w", err)
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
		return nil, fmt.Errorf("on transaction: %w", err)
	}
	return &user, nil
}

func NewService(repo Repo, ecrSvc encrypt.Service) Service {
	return &service{
		repo:   repo,
		ecrSvc: ecrSvc,
	}
}
