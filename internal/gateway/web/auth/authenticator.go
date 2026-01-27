package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
)

type Authenticator struct {
	authSvc           port.AuthService
	serviceAccountSvc port.ServiceAccountService
	apiKeyHeader      string
	userCookieName    string
}

func NewAuthenticator(apiKeyHeader, userCookieName string, authSvc port.AuthService,
	serviceAccountSvc port.ServiceAccountService,
) *Authenticator {
	return &Authenticator{
		authSvc:           authSvc,
		serviceAccountSvc: serviceAccountSvc,
		apiKeyHeader:      apiKeyHeader,
		userCookieName:    userCookieName,
	}
}

func (i *Authenticator) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rctx := ctx.Request().Context()

		if ok, err := i.authenticateByHeader(rctx, ctx.Request().Header); err != nil {
			return fmt.Errorf("authenticate by header: %w", err)
		} else if ok {
			// Identity issued by header
			return next(ctx)
		}

		if ok, err := i.authenticateByCookie(rctx, ctx.Cookies()); err != nil {
			return fmt.Errorf("authenticate by cookie: %w", err)
		} else if ok {
			// Identity issued by cookie
			return next(ctx)
		}

		// No identity issued
		return next(ctx)
	}
}

func (i *Authenticator) authenticateByHeader(ctx context.Context, header http.Header) (bool, error) {
	if ok, err := i.identityFromBearerToken(ctx, header); err != nil {
		return false, fmt.Errorf("getting identity from bearer token: %w", err)
	} else if ok {
		return true, nil
	}

	if ok, err := i.identityFromAPIKey(ctx, header); err != nil {
		return false, fmt.Errorf("getting identity from API key: %w", err)
	} else if ok {
		return true, nil
	}

	return false, nil
}

func (i *Authenticator) identityFromBearerToken(ctx context.Context, header http.Header) (bool, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return false, nil
	}

	token, ok := strings.CutPrefix(auth, "Bearer ")
	if !ok {
		return false, nil
	}

	payload, err := i.authSvc.VerifyUserToken(ctx, token)
	if err != nil {
		return false, fmt.Errorf("verifying user token: %w", err)
	}

	identity := domain.NewUserTokenIdentity(payload)
	i.registerIdentity(ctx, identity)
	return true, nil
}

func (i *Authenticator) identityFromAPIKey(ctx context.Context, header http.Header) (bool, error) {
	apiKey := header.Get(i.apiKeyHeader)
	if apiKey == "" {
		return false, nil
	}

	account, err := i.serviceAccountSvc.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return false, fmt.Errorf("getting service account by API key: %w", err)
	}

	identity := domain.NewServiceAccountIdentity(account)
	i.registerIdentity(ctx, identity)
	return true, nil
}

func (i *Authenticator) authenticateByCookie(ctx context.Context, cookies []*http.Cookie) (bool, error) {
	cookie, ok := lo.Find(cookies, func(c *http.Cookie) bool { return c.Name == i.userCookieName })
	if !ok {
		return false, nil
	}

	payload, err := i.authSvc.VerifyUserToken(ctx, cookie.Value)
	if err != nil {
		return false, fmt.Errorf("verifying user token: %w", err)
	}

	identity := domain.NewUserTokenIdentity(payload)
	i.registerIdentity(ctx, identity)
	return true, nil
}

func (i *Authenticator) registerIdentity(ctx context.Context, id domain.Identity) {
	if bag, ok := contextbag.BagFromContext(ctx); ok {
		bag.Identity = id
	}
}
