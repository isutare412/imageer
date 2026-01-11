package port

import (
	"context"
	"database/sql"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type Transactioner interface {
	BeginTx(ctx context.Context, opts ...*sql.TxOptions) (ctxWithTx context.Context, commit, rollback func() error)
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (domain.User, error)
	Upsert(ctx context.Context, user domain.User) (domain.User, error)
}

type ProjectRepository interface {
	FindByID(ctx context.Context, id string) (domain.Project, error)
	List(ctx context.Context, params domain.ListProjectsParams) (domain.Projects, error)
	Create(ctx context.Context, proj domain.Project) (domain.Project, error)
	Update(ctx context.Context, req domain.UpdateProjectRequest) (domain.Project, error)
	Delete(ctx context.Context, id string) error
}

type ServiceAccountRepository interface {
	FindByID(ctx context.Context, id string) (domain.ServiceAccount, error)
	FindByAPIKeyHash(ctx context.Context, hash string) (domain.ServiceAccount, error)
	List(ctx context.Context, params domain.ListServiceAccountsParams) (domain.ServiceAccounts, error)
	Create(ctx context.Context, sa domain.ServiceAccount) (domain.ServiceAccount, error)
	Update(ctx context.Context, req domain.UpdateServiceAccountRequest) (domain.ServiceAccount, error)
	Delete(ctx context.Context, id string) error
}

type ImageRepository interface {
	FindByID(ctx context.Context, id string) (domain.Image, error)
	Create(context.Context, domain.Image) (domain.Image, error)
	Update(context.Context, domain.UpdateImageRequest) (domain.Image, error)
}

type ImageVariantRepository interface {
	Create(context.Context, domain.ImageVariant) (domain.ImageVariant, error)
	Update(context.Context, domain.UpdateImageVariantRequest) (domain.ImageVariant, error)
}

type ImageProcessingLogRepository interface {
	Create(context.Context, domain.ImageProcessingLog) (domain.ImageProcessingLog, error)
}

type PresetRepository interface {
	FindByID(ctx context.Context, id string) (domain.Preset, error)
	FindByName(ctx context.Context, projectID, name string) (domain.Preset, error)
	List(context.Context, domain.ListPresetsParams) ([]domain.Preset, error)
}
