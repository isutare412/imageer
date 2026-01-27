package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
)

// User handlers

// GetCurrentUser gets current user details
func (h *handler) GetCurrentUser(ctx echo.Context) error {
	rctx := ctx.Request().Context()

	bag, ok := contextbag.BagFromContext(rctx)
	if !ok || bag.Identity == nil {
		return apperr.NewError(apperr.CodeUnauthorized).
			WithSummary("No authentication provided")
	}

	identity, ok := bag.Identity.(domain.UserTokenIdentity)
	if !ok {
		return apperr.NewError(apperr.CodeForbidden).
			WithSummary("Must be user token authentication")
	}

	user, err := h.userSvc.GetByID(rctx, identity.Payload.UserID)
	if err != nil {
		return fmt.Errorf("getting user by id: %w", err)
	}

	return ctx.JSON(http.StatusOK, UserToWeb(user))
}
