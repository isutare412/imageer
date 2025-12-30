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

type Service struct {
	oidcProvider port.OIDCProvider
	crypter      port.Crypter
	jwtSigner    port.JWTSigner
	jwtVerifier  port.JWTVerifier
	userRepo     port.UserRepository
	cfg          ServiceConfig
}

func NewService(cfg ServiceConfig, oidcProvider port.OIDCProvider, crypter port.Crypter,
	jwtSigner port.JWTSigner, jwtVerifier port.JWTVerifier, userRepo port.UserRepository,
) *Service {
	return &Service{
		oidcProvider: oidcProvider,
		crypter:      crypter,
		jwtSigner:    jwtSigner,
		jwtVerifier:  jwtVerifier,
		userRepo:     userRepo,
		cfg:          cfg,
	}
}

func (s *Service) StartGoogleSignIn(ctx context.Context, req domain.StartGoogleSignInRequest,
) (domain.StartGoogleSignInResponse, error) {
	state, err := s.createOIDCState(req.HTTPReq)
	if err != nil {
		return domain.StartGoogleSignInResponse{}, fmt.Errorf("creating OIDC state: %w", err)
	}

	return domain.StartGoogleSignInResponse{
		RedirectURL: s.oidcProvider.BuildAuthenticationURL(httpBaseURL(req.HTTPReq), state),
		OIDCCookie:  s.createOIDCStateCookie(state),
	}, nil
}

func (s *Service) FinishGoogleSignIn(ctx context.Context, req domain.FinishGoogleSignInRequest,
) (domain.FinishGoogleSignInResponse, error) {
	state, err := s.decryptOIDCState(req.State)
	if err != nil {
		return domain.FinishGoogleSignInResponse{}, fmt.Errorf("decrypting OIDC state: %w", err)
	}

	idToken, err := s.oidcProvider.ExchangeCode(ctx, httpBaseURL(req.HTTPReq), req.AuthCode)
	if err != nil {
		return domain.FinishGoogleSignInResponse{}, fmt.Errorf("exchanging code: %w", err)
	}

	user := domain.User{
		Role:     users.RoleGuest, // default to guest
		Nickname: idToken.FullName,
		Email:    idToken.Email,
		PhotoURL: lo.FromPtr(idToken.PictureURL),
	}

	user, err = s.userRepo.Upsert(ctx, user)
	if err != nil {
		return domain.FinishGoogleSignInResponse{}, fmt.Errorf("upserting user: %w", err)
	}

	issuedAt := time.Now()
	userPayload := domain.UserTokenPayload{
		UserID:     user.ID,
		IssuedAt:   issuedAt,
		ExpireAt:   issuedAt.Add(s.cfg.UserCookieTTL),
		Role:       user.Role,
		Nickname:   user.Nickname,
		Email:      user.Email,
		PictureURL: user.PhotoURL,
	}

	token, err := s.jwtSigner.SignUserToken(userPayload)
	if err != nil {
		return domain.FinishGoogleSignInResponse{}, fmt.Errorf("signing user token: %w", err)
	}

	return domain.FinishGoogleSignInResponse{
		RedirectURL: state.OriginURL,
		OIDCCookie:  s.deleteOIDCStateCookie(),
		UserCookie:  s.createUserCookie(token),
	}, nil
}

func (s *Service) VerifyUserToken(
	ctx context.Context, userToken string,
) (domain.UserTokenPayload, error) {
	payload, err := s.jwtVerifier.VerifyUserToken(userToken)
	if err != nil {
		return domain.UserTokenPayload{}, fmt.Errorf("verifying user token: %w", err)
	}
	return payload, nil
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
