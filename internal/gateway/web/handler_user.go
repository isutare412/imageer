package web

import (
	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

// User handlers

// GetCurrentUser gets current user details
func (h *handler) GetCurrentUser(ctx echo.Context) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}
