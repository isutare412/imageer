package auth

import (
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
)

type Authorizer struct {
	permissionInspectors []permissionInspector
	resourceInspectors   *resourceInspector
}

func NewAuthorizer(serviceAccountSvc port.ServiceAccountService, projectSvc port.ProjectService,
	imageSvc port.ImageService,
) *Authorizer {
	return &Authorizer{
		permissionInspectors: []permissionInspector{
			newAdminPermissionInspector(),
			newProjectPermissionInspector(),
		},
		resourceInspectors: newResourceInspector(serviceAccountSvc, projectSvc, imageSvc),
	}
}

func (a *Authorizer) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var identity domain.Identity
		if bag, ok := contextbag.BagFromContext(ctx); ok {
			identity = bag.Identity
		}

		// Inspect permissions
		for _, inspector := range a.permissionInspectors {
			if !inspector.isTarget(r) {
				continue
			}
			if err := inspector.inspect(r, identity); err != nil {
				gen.RespondError(w, r, fmt.Errorf("checking authorization: %w", err))
				return
			}
		}

		// Inspect existence, consistency of requested resources.
		if err := a.resourceInspectors.inspect(r); err != nil {
			gen.RespondError(w, r, fmt.Errorf("inspecting resources: %w", err))
			return
		}

		next.ServeHTTP(w, r)
	})
}
