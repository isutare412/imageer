package domain

import (
	"net/http"
	"time"

	"github.com/isutare412/imageer/pkg/users"
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

// UserTokenPayload represents the payload in a signed JWT for user
// authentication.
type UserTokenPayload struct {
	UserID     string
	IssuedAt   time.Time
	ExpireAt   time.Time
	Authority  users.Authority
	Nickname   string
	Email      string
	PictureURL string
}

type StartGoogleSignInRequest struct {
	HTTPReq *http.Request
}

type StartGoogleSignInResponse struct {
	OIDCCookie string
}
