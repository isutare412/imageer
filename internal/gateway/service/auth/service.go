package auth

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/pkg/users"
)

type AuthService struct {
	oidcProvider port.OIDCProvider
	crypter      port.Crypter
	jwtSigner    port.JWTSigner
	userRepo     port.UserRepository
	cfg          ServiceConfig
}

func NewAuthService(cfg ServiceConfig, oidcProvider port.OIDCProvider, crypter port.Crypter,
	jwtSigner port.JWTSigner, userRepo port.UserRepository) *AuthService {
	return &AuthService{
		oidcProvider: oidcProvider,
		crypter:      crypter,
		userRepo:     userRepo,
		jwtSigner:    jwtSigner,
		cfg:          cfg,
	}
}

func (s *AuthService) StartGoogleSignIn(ctx context.Context, req domain.StartGoogleSignInRequest,
) (resp domain.StartGoogleSignInResponse, err error) {
	state, err := s.createOIDCState(req.HTTPReq)
	if err != nil {
		return resp, fmt.Errorf("creating OIDC state: %w", err)
	}

	authURL := s.oidcProvider.BuildAuthenticationURL(httpBaseURL(req.HTTPReq), state)

	resp.RedirectURL = authURL
	resp.OIDCCookie = s.createOIDCStateCookie(state)
	return resp, nil
}

func (s *AuthService) FinishGoogleSignIn(ctx context.Context, req domain.FinishGoogleSignInRequest,
) (resp domain.FinishGoogleSignInResponse, err error) {
	state, err := s.decryptOIDCState(req.State)
	if err != nil {
		return resp, fmt.Errorf("decrypting OIDC state: %w", err)
	}

	idToken, err := s.oidcProvider.ExchangeCode(ctx, httpBaseURL(req.HTTPReq), req.AuthCode)
	if err != nil {
		return resp, fmt.Errorf("exchanging code: %w", err)
	}

	user := domain.User{
		Authority: users.AuthorityGuest, // default to guest
		Nickname:  idToken.FullName,
		Email:     idToken.Email,
		PhotoURL:  lo.FromPtr(idToken.PictureURL),
	}

	user, err = s.userRepo.Upsert(ctx, user)
	if err != nil {
		return resp, fmt.Errorf("upserting user: %w", err)
	}

	issuedAt := time.Now()
	userPayload := domain.UserTokenPayload{
		UserID:     user.ID,
		IssuedAt:   issuedAt,
		ExpireAt:   issuedAt.Add(s.cfg.UserCookieTTL),
		Authority:  user.Authority,
		Nickname:   user.Nickname,
		Email:      user.Email,
		PictureURL: user.PhotoURL,
	}

	token, err := s.jwtSigner.SignUserToken(userPayload)
	if err != nil {
		return resp, fmt.Errorf("signing user token: %w", err)
	}

	resp.RedirectURL = state.OriginURL
	resp.OIDCCookie = s.deleteOIDCStateCookie()
	resp.UserCookie = s.createUserCookie(token)
	return resp, nil
}

func httpBaseURL(r *http.Request) string {
	scheme := "http"
	switch {
	case r.Header.Get("X-Forwarded-Proto") == "https":
		fallthrough
	case r.TLS != nil:
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func httpReferer(r *http.Request) string {
	return cmp.Or(r.Referer(), httpBaseURL(r))
}
