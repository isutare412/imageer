package immigration

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

type permissionInspector interface {
	isTarget(r *http.Request) bool
	inspect(r *http.Request, passport domain.Passport) error
}

type adminPermissionInspector struct {
	pathPattern *regexp.Regexp
}

func newAdminPermissionInspector() *adminPermissionInspector {
	return &adminPermissionInspector{
		pathPattern: regexp.MustCompile(`^/api/v1/admin.*`),
	}
}

func (i *adminPermissionInspector) isTarget(r *http.Request) bool {
	path, err := mux.CurrentRoute(r).GetPathTemplate()
	if err != nil {
		panic(fmt.Errorf("getting path template of current request: %w", err))
	}

	return i.pathPattern.MatchString(path)
}

func (i *adminPermissionInspector) inspect(r *http.Request, passport domain.Passport) error {
	if passport == nil {
		return apperr.NewError(apperr.CodeUnauthorized).WithSummary("Need authentication")
	}

	switch pp := passport.(type) {
	case domain.UserTokenPassport:
		if !pp.Payload.IsAdmin() {
			return apperr.NewError(apperr.CodeForbidden).WithSummary("Admin role required")
		}

	case domain.ServiceAccountPassport:
		if !pp.ServiceAccount.HasFullAccess() {
			return apperr.NewError(apperr.CodeForbidden).WithSummary("Full access scope required")
		}

	default:
		return apperr.NewError(apperr.CodeForbidden).WithSummary("Unexpected passport type")
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

func (i *projectPermissionInspector) isTarget(r *http.Request) bool {
	path, err := mux.CurrentRoute(r).GetPathTemplate()
	if err != nil {
		panic(fmt.Errorf("getting path template of current request: %w", err))
	}

	return i.pathPattern.MatchString(path)
}

func (i *projectPermissionInspector) inspect(r *http.Request, passport domain.Passport) error {
	if passport == nil {
		return apperr.NewError(apperr.CodeUnauthorized).WithSummary("Need authentication")
	}

	switch pp := passport.(type) {
	case domain.UserTokenPassport:
		if !pp.Payload.IsAdmin() {
			return apperr.NewError(apperr.CodeForbidden).WithSummary("Admin role required")
		}

	case domain.ServiceAccountPassport:
		switch pp.ServiceAccount.AccessScope {
		case serviceaccounts.AccessScopeFull:
			// Has full access, allow.
			return nil

		case serviceaccounts.AccessScopeProject:
			projID := mux.Vars(r)["projectId"]
			if projID == "" {
				return apperr.NewError(apperr.CodeBadRequest).WithSummary("Project ID required")
			}

			_, projFound := lo.Find(pp.ServiceAccount.Projects,
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
		return apperr.NewError(apperr.CodeForbidden).WithSummary("Unexpected passport type")
	}

	return nil
}
