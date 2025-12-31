package immigration

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
)

type Immigration struct {
	permissionInspectors []permissionInspector
	resourceInspectors   *resourceInspector
}

func New(serviceAccountSvc port.ServiceAccountService) *Immigration {
	return &Immigration{
		permissionInspectors: []permissionInspector{
			newAdminPermissionInspector(),
			newProjectPermissionInspector(),
		},
		resourceInspectors: newResourceInspector(serviceAccountSvc),
	}
}

func (i *Immigration) Immigrate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rctx := ctx.Request().Context()

		var passport domain.Passport
		if bag, ok := contextbag.BagFromContext(rctx); ok {
			passport = bag.Passport
		}

		// Inspect permissions
		for _, inspector := range i.permissionInspectors {
			if !inspector.isTarget(ctx) {
				continue
			}
			if err := inspector.inspect(ctx, passport); err != nil {
				return fmt.Errorf("inspecting passport: %w", err)
			}
		}

		// Inspect existence, consistency of requested resources.
		if err := i.resourceInspectors.inspect(ctx); err != nil {
			return fmt.Errorf("inspecting resources: %w", err)
		}

		return next(ctx)
	}
}
