package immigration

import (
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
)

type Immigration struct {
	permissionInspectors []permissionInspector
	resourceInspectors   *resourceInspector
}

func New(serviceAccountSvc port.ServiceAccountService, projectSvc port.ProjectService,
	imageSvc port.ImageService,
) *Immigration {
	return &Immigration{
		permissionInspectors: []permissionInspector{
			newAdminPermissionInspector(),
			newProjectPermissionInspector(),
		},
		resourceInspectors: newResourceInspector(serviceAccountSvc, projectSvc, imageSvc),
	}
}

func (i *Immigration) Immigrate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var passport domain.Passport
		if bag, ok := contextbag.BagFromContext(ctx); ok {
			passport = bag.Passport
		}

		// Inspect permissions
		for _, inspector := range i.permissionInspectors {
			if !inspector.isTarget(r) {
				continue
			}
			if err := inspector.inspect(r, passport); err != nil {
				gen.RespondError(w, r, fmt.Errorf("inspecting passport: %w", err))
				return
			}
		}

		// Inspect existence, consistency of requested resources.
		if err := i.resourceInspectors.inspect(r); err != nil {
			gen.RespondError(w, r, fmt.Errorf("inspecting resources: %w", err))
			return
		}

		next.ServeHTTP(w, r)
	})
}
