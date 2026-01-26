package webv2

import (
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/apperr"
)

// User handlers

// GetCurrentUser gets current user details
func (h *handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bag, ok := contextbag.BagFromContext(ctx)
	if !ok || bag.Passport == nil {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeUnauthorized).
			WithSummary("No authentication provided"))
		return
	}

	passport, ok := bag.Passport.(domain.UserTokenPassport)
	if !ok {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeForbidden).
			WithSummary("Must be user token authentication"))
		return
	}

	user, err := h.userSvc.GetByID(ctx, passport.Payload.UserID)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting user by id: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, UserToWeb(user))
}
