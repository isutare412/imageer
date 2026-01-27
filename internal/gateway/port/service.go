package port

import (
	"context"

	"github.com/isutare412/imageer/internal/gateway/domain"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type AuthService interface {
	StartGoogleSignIn(context.Context, domain.StartGoogleSignInRequest) (domain.StartGoogleSignInResponse, error)
	FinishGoogleSignIn(context.Context, domain.FinishGoogleSignInRequest) (domain.FinishGoogleSignInResponse, error)
	SignOut(ctx context.Context) domain.SignOutResponse
	VerifyUserToken(ctx context.Context, userToken string) (domain.UserTokenPayload, error)
	RefreshUserToken(ctx context.Context, userID string) (domain.RefreshUserTokenResponse, error)
}

type UserService interface {
	GetByID(ctx context.Context, id string) (domain.User, error)
}

type ServiceAccountService interface {
	GetByID(ctx context.Context, id string) (domain.ServiceAccount, error)
	GetByAPIKey(ctx context.Context, key string) (domain.ServiceAccount, error)
	List(context.Context, domain.ListServiceAccountsParams) (domain.ServiceAccounts, error)
	Create(context.Context, domain.CreateServiceAccountRequest) (domain.ServiceAccountWithAPIKey, error)
	Update(context.Context, domain.UpdateServiceAccountRequest) (domain.ServiceAccount, error)
	Delete(ctx context.Context, id string) error
}

type ProjectService interface {
	GetByID(ctx context.Context, id string) (domain.Project, error)
	List(context.Context, domain.ListProjectsParams) (domain.Projects, error)
	Create(context.Context, domain.CreateProjectRequest) (domain.Project, error)
	Update(context.Context, domain.UpdateProjectRequest) (domain.Project, error)
	Delete(ctx context.Context, id string) error
}

type ImageService interface {
	Get(ctx context.Context, imageID string) (domain.Image, error)
	GetWaitUntilProcessed(ctx context.Context, imageID string) (domain.Image, error)
	List(context.Context, domain.ListImagesParams) (domain.Images, error)
	Delete(ctx context.Context, id string) error
	DeleteS3Objects(context.Context, *imageerv1.ImageS3DeleteRequest) error
	CreateUploadURL(context.Context, domain.CreateUploadURLRequest) (domain.UploadURL, error)
	StartImageProcessingOnUpload(ctx context.Context, s3Key string) error
	ReceiveImageProcessResult(context.Context, *imageerv1.ImageProcessResult) error
}
