package handlers

import "github.com/isutare412/imageer/internal/gateway/port"

// handler implements the ServerInterface for handling HTTP requests
type handler struct {
	authSvc           port.AuthService
	serviceAccountSvc port.ServiceAccountService
	projectSvc        port.ProjectService
	userSvc           port.UserService
	imageSvc          port.ImageService
}

// NewHandler creates a new Handler instance
func NewHandler(authSvc port.AuthService, serviceAccountSvc port.ServiceAccountService,
	projectSvc port.ProjectService, userSvc port.UserService, imageSvc port.ImageService,
) *handler {
	return &handler{
		authSvc:           authSvc,
		serviceAccountSvc: serviceAccountSvc,
		projectSvc:        projectSvc,
		userSvc:           userSvc,
		imageSvc:          imageSvc,
	}
}
