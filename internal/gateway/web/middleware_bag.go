package web

import (
	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
)

func withContextBag(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rctx := ctx.Request().Context()
		rctx = contextbag.WithBag(rctx)
		ctx.SetRequest(ctx.Request().WithContext(rctx))
		return next(ctx)
	}
}
