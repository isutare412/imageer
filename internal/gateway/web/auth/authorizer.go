package auth

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
)

type Authorizer struct {
	permissionInspectors []permissionInspector
	resourceInspectors   *resourceInspector
}

func NewAuthorizer(serviceAccountSvc port.ServiceAccountService, projectSvc port.ProjectService,
	imageSvc port.ImageService,
) *Authorizer {
	return &Authorizer{
		permissionInspectors: []permissionInspector{
			newAdminPermissionInspector(),
			newProjectPermissionInspector(),
		},
		resourceInspectors: newResourceInspector(serviceAccountSvc, projectSvc, imageSvc),
	}
}

func (i *Authorizer) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rctx := ctx.Request().Context()

		var identity domain.Identity
		if bag, ok := contextbag.BagFromContext(rctx); ok {
			identity = bag.Identity
		}

		// Inspect permissions
		for _, inspector := range i.permissionInspectors {
			if !inspector.isTarget(ctx) {
				continue
			}
			if err := inspector.inspect(ctx, identity); err != nil {
				return fmt.Errorf("checking authorization: %w", err)
			}
		}

		// Inspect existence, consistency of requested resources.
		if err := i.resourceInspectors.inspect(ctx); err != nil {
			return fmt.Errorf("inspecting resources: %w", err)
		}

		return next(ctx)
	}
}
