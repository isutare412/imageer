package auth

import (
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

type permissionInspector interface {
	isTarget(ctx echo.Context) bool
	inspect(ctx echo.Context, identity domain.Identity) error
}

type adminPermissionInspector struct {
	pathPattern *regexp.Regexp
}

func newAdminPermissionInspector() *adminPermissionInspector {
	return &adminPermissionInspector{
		pathPattern: regexp.MustCompile(`^/api/v1/admin.*`),
	}
}

func (i *adminPermissionInspector) isTarget(ctx echo.Context) bool {
	return i.pathPattern.MatchString(ctx.Path())
}

func (i *adminPermissionInspector) inspect(ctx echo.Context, identity domain.Identity) error {
	if identity == nil {
		return apperr.NewError(apperr.CodeUnauthorized).WithSummary("Need authentication")
	}

	switch id := identity.(type) {
	case domain.UserTokenIdentity:
		if !id.Payload.IsAdmin() {
			return apperr.NewError(apperr.CodeForbidden).WithSummary("Admin role required")
		}

	case domain.ServiceAccountIdentity:
		if !id.ServiceAccount.HasFullAccess() {
			return apperr.NewError(apperr.CodeForbidden).WithSummary("Full access scope required")
		}

	default:
		return apperr.NewError(apperr.CodeForbidden).WithSummary("Unexpected identity type")
	}

	return nil
}

type projectPermissionInspector struct {
	pathPattern *regexp.Regexp
}

func newProjectPermissionInspector() *projectPermissionInspector {
	return &projectPermissionInspector{
		pathPattern: regexp.MustCompile(`^/api/v1/projects/([^/]+).*`),
	}
}

func (i *projectPermissionInspector) isTarget(ctx echo.Context) bool {
	return i.pathPattern.MatchString(ctx.Path())
}

func (i *projectPermissionInspector) inspect(ctx echo.Context, identity domain.Identity) error {
	if identity == nil {
		return apperr.NewError(apperr.CodeUnauthorized).WithSummary("Need authentication")
	}

	switch id := identity.(type) {
	case domain.UserTokenIdentity:
		if !id.Payload.IsAdmin() {
			return apperr.NewError(apperr.CodeForbidden).WithSummary("Admin role required")
		}

	case domain.ServiceAccountIdentity:
		switch id.ServiceAccount.AccessScope {
		case serviceaccounts.AccessScopeFull:
			// Has full access, allow.
			return nil

		case serviceaccounts.AccessScopeProject:
			projID := ctx.Param("projectId")
			if projID == "" {
				return apperr.NewError(apperr.CodeBadRequest).WithSummary("Project ID required")
			}

			_, projFound := lo.Find(id.ServiceAccount.Projects,
				func(p domain.ProjectReference) bool {
					return p.ID == projID
				})
			if !projFound {
				return apperr.NewError(apperr.CodeForbidden).WithSummary("Insufficient service account access scope")
			}

		default:
			return apperr.NewError(apperr.CodeForbidden).WithSummary("Insufficient service account access scope")
		}

	default:
		return apperr.NewError(apperr.CodeForbidden).WithSummary("Unexpected identity type")
	}

	return nil
}
