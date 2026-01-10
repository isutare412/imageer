package port

import (
	"context"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type AuthService interface {
	StartGoogleSignIn(context.Context, domain.StartGoogleSignInRequest) (domain.StartGoogleSignInResponse, error)
	FinishGoogleSignIn(context.Context, domain.FinishGoogleSignInRequest) (domain.FinishGoogleSignInResponse, error)
	VerifyUserToken(ctx context.Context, userToken string) (domain.UserTokenPayload, error)
}

type ServiceAccountService interface {
	GetByID(ctx context.Context, id string) (domain.ServiceAccount, error)
	GetByAPIKey(ctx context.Context, key string) (domain.ServiceAccount, error)
	List(ctx context.Context, params domain.ListServiceAccountsParams) (domain.ServiceAccounts, error)
	Create(ctx context.Context, req domain.CreateServiceAccountRequest) (domain.ServiceAccountWithAPIKey, error)
	Update(ctx context.Context, req domain.UpdateServiceAccountRequest) (domain.ServiceAccount, error)
	Delete(ctx context.Context, id string) error
}

type ProjectService interface {
	GetByID(ctx context.Context, id string) (domain.Project, error)
	List(ctx context.Context, params domain.ListProjectsParams) (domain.Projects, error)
	Create(ctx context.Context, req domain.CreateProjectRequest) (domain.Project, error)
	Update(ctx context.Context, req domain.UpdateProjectRequest) (domain.Project, error)
	Delete(ctx context.Context, id string) error
}

type ImageService interface {
	CreateUploadURL(context.Context, domain.CreateUploadURLRequest) (domain.UploadURL, error)
	StartImageProcessingOnUpload(ctx context.Context, s3Key string) error
}
