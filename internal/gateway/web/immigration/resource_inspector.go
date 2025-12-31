package immigration

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
)

type resourceInspector struct {
	serviceAccountSvc port.ServiceAccountService
	projectSvc        port.ProjectService
}

func newResourceInspector(serviceAccountSvc port.ServiceAccountService,
	projectSvc port.ProjectService,
) *resourceInspector {
	return &resourceInspector{
		serviceAccountSvc: serviceAccountSvc,
		projectSvc:        projectSvc,
	}
}

func (i *resourceInspector) inspect(ctx echo.Context) error {
	rctx := ctx.Request().Context()

	if _, _, err := i.fetchServiceAccount(rctx, ctx); err != nil {
		return fmt.Errorf("fetching requested service account: %w", err)
	}

	if _, _, err := i.fetchProject(rctx, ctx); err != nil {
		return fmt.Errorf("fetching requested project: %w", err)
	}

	return nil
}

func (i *resourceInspector) fetchServiceAccount(ctx context.Context, ectx echo.Context,
) (domain.ServiceAccount, bool, error) {
	id := ectx.Param("serviceAccountId")
	if id == "" {
		return domain.ServiceAccount{}, false, nil
	}

	account, err := i.serviceAccountSvc.GetByID(ctx, id)
	if err != nil {
		return domain.ServiceAccount{}, false, fmt.Errorf("getting service account by id: %w", err)
	}
	return account, true, nil
}

func (i *resourceInspector) fetchProject(ctx context.Context, ectx echo.Context,
) (domain.Project, bool, error) {
	id := ectx.Param("projectId")
	if id == "" {
		return domain.Project{}, false, nil
	}

	project, err := i.projectSvc.GetByID(ctx, id)
	if err != nil {
		return domain.Project{}, false, fmt.Errorf("getting project by id: %w", err)
	}
	return project, true, nil
}
