package user

import (
	"context"
	"fmt"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
)

type Service struct {
	userRepo port.UserRepository
}

func NewService(userRepo port.UserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("finding user by id: %w", err)
	}
	return user, nil
}
