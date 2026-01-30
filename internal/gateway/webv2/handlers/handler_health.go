package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
)

func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	gen.RespondNoContent(w, http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eg, egCtx := errgroup.WithContext(ctx)
	for _, checker := range h.healthCheckers {
		eg.Go(func() error {
			if err := checker.HealthCheck(egCtx); err != nil {
				return fmt.Errorf("health checking %s: %w", checker.ComponentName(), err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		gen.RespondError(w, r, err)
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}
