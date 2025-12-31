package web

import "github.com/isutare412/imageer/internal/gateway/port"

// handler implements the ServerInterface for handling HTTP requests
type handler struct {
	authSvc           port.AuthService
	serviceAccountSvc port.ServiceAccountService
	projectSvc        port.ProjectService
}

// newHandler creates a new Handler instance
func newHandler(authSvc port.AuthService, serviceAccountSvc port.ServiceAccountService,
	projectSvc port.ProjectService,
) *handler {
	return &handler{
		authSvc:           authSvc,
		serviceAccountSvc: serviceAccountSvc,
		projectSvc:        projectSvc,
	}
}
