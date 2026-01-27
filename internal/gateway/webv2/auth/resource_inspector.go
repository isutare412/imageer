package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/pkg/apperr"
)

type resourceInspector struct {
	serviceAccountSvc port.ServiceAccountService
	projectSvc        port.ProjectService
	imageSvc          port.ImageService
}

func newResourceInspector(serviceAccountSvc port.ServiceAccountService,
	projectSvc port.ProjectService, imageSvc port.ImageService,
) *resourceInspector {
	return &resourceInspector{
		serviceAccountSvc: serviceAccountSvc,
		projectSvc:        projectSvc,
		imageSvc:          imageSvc,
	}
}

func (i *resourceInspector) inspect(r *http.Request) error {
	ctx := r.Context()

	if _, _, err := i.fetchServiceAccount(ctx, r); err != nil {
		return fmt.Errorf("fetching requested service account: %w", err)
	}

	project, projectExists, err := i.fetchProject(ctx, r)
	if err != nil {
		return fmt.Errorf("fetching requested project: %w", err)
	}

	image, imageExists, err := i.fetchImage(ctx, r)
	if err != nil {
		return fmt.Errorf("fetching requested image: %w", err)
	}

	if imageExists && projectExists {
		if image.Project.ID != project.ID {
			return apperr.NewError(apperr.CodeNotFound).
				WithSummary("Image not found in the specified project")
		}
	}

	return nil
}

func (i *resourceInspector) fetchServiceAccount(ctx context.Context, r *http.Request,
) (domain.ServiceAccount, bool, error) {
	id := mux.Vars(r)["serviceAccountId"]
	if id == "" {
		return domain.ServiceAccount{}, false, nil
	}

	account, err := i.serviceAccountSvc.GetByID(ctx, id)
	if err != nil {
		return domain.ServiceAccount{}, false, fmt.Errorf("getting service account by id: %w", err)
	}
	return account, true, nil
}

func (i *resourceInspector) fetchProject(ctx context.Context, r *http.Request,
) (domain.Project, bool, error) {
	id := mux.Vars(r)["projectId"]
	if id == "" {
		return domain.Project{}, false, nil
	}

	project, err := i.projectSvc.GetByID(ctx, id)
	if err != nil {
		return domain.Project{}, false, fmt.Errorf("getting project by id: %w", err)
	}
	return project, true, nil
}

func (i *resourceInspector) fetchImage(ctx context.Context, r *http.Request,
) (domain.Image, bool, error) {
	id := mux.Vars(r)["imageId"]
	if id == "" {
		return domain.Image{}, false, nil
	}

	image, err := i.imageSvc.Get(ctx, id)
	if err != nil {
		return domain.Image{}, false, fmt.Errorf("getting image by id: %w", err)
	}
	return image, true, nil
}
