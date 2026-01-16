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
	Role       users.Role
	Nickname   string
	Email      string
	PictureURL string
}

func (p UserTokenPayload) IsAdmin() bool {
	return p.Role == users.RoleAdmin
}

type OIDCState struct {
	OriginURL string `json:"originUrl"`
}

type StartGoogleSignInRequest struct {
	HTTPReq *http.Request
}

type StartGoogleSignInResponse struct {
	RedirectURL string
	OIDCCookie  *http.Cookie
}

type FinishGoogleSignInRequest struct {
	HTTPReq  *http.Request
	AuthCode string
	State    string
}

type FinishGoogleSignInResponse struct {
	RedirectURL string
	OIDCCookie  *http.Cookie
	UserCookie  *http.Cookie
}

type SignOutResponse struct {
	UserCookie *http.Cookie
}
