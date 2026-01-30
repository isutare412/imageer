package handlers

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
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bag, ok := contextbag.BagFromContext(ctx)
	if !ok || bag.Identity == nil {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeUnauthorized).
			WithSummary("No authentication provided"))
		return
	}

	identity, ok := bag.Identity.(domain.UserTokenIdentity)
	if !ok {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeForbidden).
			WithSummary("Must be user token authentication"))
		return
	}

	user, err := h.userSvc.GetByID(ctx, identity.Payload.UserID)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting user by id: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, UserToWeb(user))
}
