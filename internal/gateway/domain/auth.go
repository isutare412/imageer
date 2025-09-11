package domain

import (
	"net/http"
	"time"
)

// IDTokenPayload represents the payload extracted from an ID token by IdP.
// Pointer fields are optional.
type IDTokenPayload struct {
	Audience      string
	Issuer        string
	Subject       string
	Email         string
	EmailVerified *bool
	FamilyName    *string
	GivenName     *string
	FullName      string
	PictureURL    *string
	ProfileURL    *string
	Expiry        time.Time
	IssuedAt      time.Time
}

type StartGoogleSignInRequest struct {
	HTTPReq *http.Request
}

type StartGoogleSignInResponse struct {
	OIDCCookie string
}
