package immigration

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
)

type Immigration struct {
	inspectors []inspector
}

func New() *Immigration {
	return &Immigration{
		inspectors: []inspector{
			newAdminInspector(),
			newProjectPermissionInspector(),
		},
	}
}

func (i *Immigration) Immigrate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rctx := ctx.Request().Context()

		var passport domain.Passport
		if bag, ok := contextbag.BagFromContext(rctx); ok {
			passport = bag.Passport
		}

		for _, inspector := range i.inspectors {
			if !inspector.isTarget(ctx) {
				continue
			}

			if err := inspector.inspect(ctx, passport); err != nil {
				return fmt.Errorf("inspecting passport: %w", err)
			}
		}

		return next(ctx)
	}
}
