package web

import (
	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

// Authentication handlers

// StartGoogleSignIn starts Google Sign-In process
func (h *handler) StartGoogleSignIn(ctx echo.Context) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// FinishGoogleSignIn finishes Google Sign-In process
func (h *handler) FinishGoogleSignIn(ctx echo.Context, params FinishGoogleSignInParams) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}
