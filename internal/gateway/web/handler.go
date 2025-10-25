package web

import "github.com/isutare412/imageer/internal/gateway/port"

// handler implements the ServerInterface for handling HTTP requests
type handler struct {
	authSvc port.AuthService
}

// newHandler creates a new Handler instance
func newHandler(authSvc port.AuthService) *handler {
	return &handler{
		authSvc: authSvc,
	}
}
