package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
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

func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if ok, err := a.authenticateByHeader(ctx, r.Header); err != nil {
			gen.RespondError(w, r, fmt.Errorf("authenticate by header: %w", err))
			return
		} else if ok {
			// Identity issued by header
			next.ServeHTTP(w, r)
			return
		}

		if ok, err := a.authenticateByCookie(ctx, r.Cookies()); err != nil {
			gen.RespondError(w, r, fmt.Errorf("authenticate by cookie: %w", err))
			return
		} else if ok {
			// Identity issued by cookie
			next.ServeHTTP(w, r)
			return
		}

		// No identity issued
		next.ServeHTTP(w, r)
	})
}

func (a *Authenticator) authenticateByHeader(ctx context.Context, header http.Header) (bool, error) {
	if ok, err := a.identityFromBearerToken(ctx, header); err != nil {
		return false, fmt.Errorf("getting identity from bearer token: %w", err)
	} else if ok {
		return true, nil
	}

	if ok, err := a.identityFromAPIKey(ctx, header); err != nil {
		return false, fmt.Errorf("getting identity from API key: %w", err)
	} else if ok {
		return true, nil
	}

	return false, nil
}

func (a *Authenticator) identityFromBearerToken(ctx context.Context, header http.Header) (bool, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return false, nil
	}

	token, ok := strings.CutPrefix(auth, "Bearer ")
	if !ok {
		return false, nil
	}

	payload, err := a.authSvc.VerifyUserToken(ctx, token)
	if err != nil {
		return false, fmt.Errorf("verifying user token: %w", err)
	}

	identity := domain.NewUserTokenIdentity(payload)
	a.registerIdentity(ctx, identity)
	return true, nil
}

func (a *Authenticator) identityFromAPIKey(ctx context.Context, header http.Header) (bool, error) {
	apiKey := header.Get(a.apiKeyHeader)
	if apiKey == "" {
		return false, nil
	}

	account, err := a.serviceAccountSvc.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return false, fmt.Errorf("getting service account by API key: %w", err)
	}

	identity := domain.NewServiceAccountIdentity(account)
	a.registerIdentity(ctx, identity)
	return true, nil
}

func (a *Authenticator) authenticateByCookie(ctx context.Context, cookies []*http.Cookie) (bool, error) {
	cookie, ok := lo.Find(cookies, func(c *http.Cookie) bool { return c.Name == a.userCookieName })
	if !ok {
		return false, nil
	}

	payload, err := a.authSvc.VerifyUserToken(ctx, cookie.Value)
	if err != nil {
		return false, fmt.Errorf("verifying user token: %w", err)
	}

	identity := domain.NewUserTokenIdentity(payload)
	a.registerIdentity(ctx, identity)
	return true, nil
}

func (a *Authenticator) registerIdentity(ctx context.Context, id domain.Identity) {
	if bag, ok := contextbag.BagFromContext(ctx); ok {
		bag.Identity = id
	}
}
