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
}

func newResourceInspector(serviceAccountSvc port.ServiceAccountService) *resourceInspector {
	return &resourceInspector{
		serviceAccountSvc: serviceAccountSvc,
	}
}

func (i *resourceInspector) inspect(ctx echo.Context) error {
	rctx := ctx.Request().Context()

	if _, _, err := i.fetchServiceAccount(rctx, ctx); err != nil {
		return fmt.Errorf("fetching requested service account: %w", err)
	}

	return nil
}

func (i *resourceInspector) fetchServiceAccount(
	ctx context.Context, ectx echo.Context,
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
