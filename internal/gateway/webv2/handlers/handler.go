package handlers

import "github.com/isutare412/imageer/internal/gateway/port"

// Handler implements the ServerInterface for handling HTTP requests
type Handler struct {
	authSvc           port.AuthService
	serviceAccountSvc port.ServiceAccountService
	projectSvc        port.ProjectService
	userSvc           port.UserService
	imageSvc          port.ImageService
	healthCheckers    []port.HealthChecker
}

// NewHandler creates a new Handler instance
func NewHandler(
	healthCheckers []port.HealthChecker,
	authSvc port.AuthService,
	serviceAccountSvc port.ServiceAccountService,
	projectSvc port.ProjectService,
	userSvc port.UserService,
	imageSvc port.ImageService,
) *Handler {
	return &Handler{
		authSvc:           authSvc,
		serviceAccountSvc: serviceAccountSvc,
		projectSvc:        projectSvc,
		userSvc:           userSvc,
		imageSvc:          imageSvc,
		healthCheckers:    healthCheckers,
	}
}
