package port

import (
	"context"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type UserRepository interface {
	FindByID(ctx context.Context, id string) (domain.User, error)
	Upsert(ctx context.Context, user domain.User) (domain.User, error)
}

type ProjectRepository interface {
	FindByID(ctx context.Context, id string) (domain.Project, error)
	List(ctx context.Context, params domain.ListProjectsParams) ([]domain.Project, error)
	Create(ctx context.Context, proj domain.Project) (domain.Project, error)
	Update(ctx context.Context, req domain.UpdateProjectRequest) (domain.Project, error)
}

type ServiceAccountRepository interface {
	FindByID(ctx context.Context, id string) (domain.ServiceAccount, error)
	List(ctx context.Context, params domain.ListServiceAccountsParams) ([]domain.ServiceAccount, error)
	Create(ctx context.Context, sa domain.ServiceAccount) (domain.ServiceAccount, error)
	Update(ctx context.Context, req domain.UpdateServiceAccountRequest) (domain.ServiceAccount, error)
	Delete(ctx context.Context, id string) error
}
