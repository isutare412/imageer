package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

// Authentication handlers

// StartGoogleSignIn starts Google Sign-In process
func (h *handler) StartGoogleSignIn(ctx echo.Context) error {
	rctx := ctx.Request().Context()

	req := domain.StartGoogleSignInRequest{
		HTTPReq: ctx.Request(),
	}
	resp, err := h.authSvc.StartGoogleSignIn(rctx, req)
	if err != nil {
		return fmt.Errorf("start google sign-in: %w", err)
	}

	ctx.SetCookie(resp.OIDCCookie)
	return ctx.Redirect(http.StatusFound, resp.RedirectURL)
}

// FinishGoogleSignIn finishes Google Sign-In process
func (h *handler) FinishGoogleSignIn(ctx echo.Context, params FinishGoogleSignInParams) error {
	rctx := ctx.Request().Context()

	req := domain.FinishGoogleSignInRequest{
		HTTPReq:  ctx.Request(),
		AuthCode: params.Code,
		State:    params.State,
	}
	resp, err := h.authSvc.FinishGoogleSignIn(rctx, req)
	if err != nil {
		return fmt.Errorf("finish google sign-in: %w", err)
	}

	ctx.SetCookie(resp.OIDCCookie)
	ctx.SetCookie(resp.UserCookie)
	return ctx.Redirect(http.StatusFound, resp.RedirectURL)
}
