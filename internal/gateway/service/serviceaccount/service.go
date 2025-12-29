package serviceaccount

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/service/serviceaccount/apikey"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/validation"
)

type Service struct {
	serviceAccountRepo port.ServiceAccountRepository
}

func NewService(serviceAccountRepo port.ServiceAccountRepository) *Service {
	return &Service{
		serviceAccountRepo: serviceAccountRepo,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.ServiceAccount, error) {
	account, err := s.serviceAccountRepo.FindByID(ctx, id)
	if err != nil {
		return domain.ServiceAccount{}, fmt.Errorf("finding service account: %w", err)
	}
	return account, nil
}

func (s *Service) GetByAPIKey(ctx context.Context, key string) (domain.ServiceAccount, error) {
	apiKey, err := apikey.ParseString(key)
	if err != nil {
		return domain.ServiceAccount{}, fmt.Errorf("parsing API key: %w", err)
	}

	account, err := s.serviceAccountRepo.FindByAPIKeyHash(ctx, apiKey.Hash())
	switch {
	case apperr.IsErrorStatusCode(err, http.StatusNotFound):
		return domain.ServiceAccount{},
			apperr.NewError(apperr.CodeBadRequest).WithSummary("API key is invalid")
	case err != nil:
		return domain.ServiceAccount{}, fmt.Errorf("finding service account: %w", err)
	}

	return account, nil
}

func (s *Service) List(
	ctx context.Context, params domain.ListServiceAccountsParams,
) ([]domain.ServiceAccount, error) {
	if err := validation.Validate(params); err != nil {
		return nil, fmt.Errorf("validating params: %w", err)
	}

	accounts, err := s.serviceAccountRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing service accounts: %w", err)
	}

	return accounts, nil
}

func (s *Service) Create(
	ctx context.Context, req domain.CreateServiceAccountRequest,
) (domain.ServiceAccountWithAPIKey, error) {
	if err := validation.Validate(req); err != nil {
		return domain.ServiceAccountWithAPIKey{}, fmt.Errorf("validating request: %w", err)
	}

	apiKey := apikey.New()
	account := req.ToServiceAccount(apiKey.Hash())
	account, err := s.serviceAccountRepo.Create(ctx, account)
	if err != nil {
		return domain.ServiceAccountWithAPIKey{}, fmt.Errorf("creating service account: %w", err)
	}

	return domain.ServiceAccountWithAPIKey{
		ServiceAccount: account,
		APIKey:         apiKey.String(),
	}, nil
}

func (s *Service) Update(
	ctx context.Context, req domain.UpdateServiceAccountRequest,
) (domain.ServiceAccount, error) {
	if err := validation.Validate(req); err != nil {
		return domain.ServiceAccount{}, fmt.Errorf("validating request: %w", err)
	}

	account, err := s.serviceAccountRepo.Update(ctx, req)
	if err != nil {
		return domain.ServiceAccount{}, fmt.Errorf("updating service account: %w", err)
	}

	return account, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.serviceAccountRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting service account: %w", err)
	}
	return nil
}
