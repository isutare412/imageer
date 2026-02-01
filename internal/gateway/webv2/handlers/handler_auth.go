package handlers

import (
	"fmt"
	"net/http"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/trace"
)

// Authentication handlers

// StartGoogleSignIn starts Google Sign-In process
func (h *Handler) StartGoogleSignIn(w http.ResponseWriter, r *http.Request,
	params gen.StartGoogleSignInParams,
) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.StartGoogleSignIn")
	defer span.End()

	req := domain.StartGoogleSignInRequest{
		HTTPReq:      r,
		RedirectPath: lo.FromPtr(params.Redirect),
	}
	resp, err := h.authSvc.StartGoogleSignIn(ctx, req)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("start google sign-in: %w", err))
		return
	}

	http.SetCookie(w, resp.OIDCCookie)
	http.Redirect(w, r, resp.RedirectURL, http.StatusFound)
}

// FinishGoogleSignIn finishes Google Sign-In process
func (h *Handler) FinishGoogleSignIn(w http.ResponseWriter, r *http.Request,
	params gen.FinishGoogleSignInParams,
) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.FinishGoogleSignIn")
	defer span.End()

	req := domain.FinishGoogleSignInRequest{
		HTTPReq:  r,
		AuthCode: params.Code,
		State:    params.State,
	}
	resp, err := h.authSvc.FinishGoogleSignIn(ctx, req)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("finish google sign-in: %w", err))
		return
	}

	http.SetCookie(w, resp.OIDCCookie)
	http.SetCookie(w, resp.UserCookie)
	http.Redirect(w, r, resp.RedirectURL, http.StatusFound)
}

// SignOut signs out the current user by clearing the user cookie
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.SignOut")
	defer span.End()

	resp := h.authSvc.SignOut(ctx)

	http.SetCookie(w, resp.UserCookie)
	gen.RespondNoContent(w, http.StatusOK)
}
