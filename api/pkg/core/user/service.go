package user

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, user *User) (int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	UpdateCredit(ctx context.Context, id int64, delta int64) (*User, error)
}

type service struct {
	repo Repo
}

func (s *service) Create(ctx context.Context, user *User) (int64, error) {
	var newID int64
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("on create user: %w", err)
		}

		newID = user.ID
		return nil
	}); err != nil {
		return 0, fmt.Errorf("on transaction: %w", err)
	}
	return newID, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := s.repo.Session(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return fmt.Errorf("on first user: %w", err)
		}
		return nil
	}); err != nil {
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

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}
