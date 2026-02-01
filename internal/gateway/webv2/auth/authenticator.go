package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/tracing"
)

type Authenticator struct {
	authSvc               port.AuthService
	serviceAccountSvc     port.ServiceAccountService
	apiKeyHeader          string
	userCookieName        string
	tokenRefreshThreshold time.Duration
}

func NewAuthenticator(apiKeyHeader, userCookieName string, tokenRefreshThreshold time.Duration,
	authSvc port.AuthService, serviceAccountSvc port.ServiceAccountService,
) *Authenticator {
	return &Authenticator{
		authSvc:               authSvc,
		serviceAccountSvc:     serviceAccountSvc,
		apiKeyHeader:          apiKeyHeader,
		userCookieName:        userCookieName,
		tokenRefreshThreshold: tokenRefreshThreshold,
	}
}

func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracing.StartSpan(r.Context(), "webv2.Authenticator.Authenticate")
		defer span.End()

		if ok, err := a.authenticateByHeader(ctx, r.Header); err != nil {
			gen.RespondError(w, r, fmt.Errorf("authenticate by header: %w", err))
			return
		} else if ok {
			// Identity issued by header
			next.ServeHTTP(w, r)
			return
		}

		if ok, err := a.authenticateByCookie(ctx, w, r.Cookies()); err != nil {
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

func (a *Authenticator) authenticateByCookie(ctx context.Context, w http.ResponseWriter,
	cookies []*http.Cookie,
) (bool, error) {
	cookie, ok := lo.Find(cookies, func(c *http.Cookie) bool { return c.Name == a.userCookieName })
	if !ok {
		return false, nil
	}

	payload, err := a.authSvc.VerifyUserToken(ctx, cookie.Value)
	if err != nil {
		return false, fmt.Errorf("verifying user token: %w", err)
	}

	// Refresh token if near expiration
	if time.Until(payload.ExpireAt) < a.tokenRefreshThreshold {
		a.refreshUserToken(ctx, w, payload.UserID)
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

func (a *Authenticator) refreshUserToken(ctx context.Context, w http.ResponseWriter, userID string,
) {
	resp, err := a.authSvc.RefreshUserToken(ctx, userID)
	if err != nil {
		a.logRefreshError(ctx, err)
		return
	}

	slog.InfoContext(ctx, "Refreshed user token", "userId", userID)
	http.SetCookie(w, resp.UserCookie)
}

func (a *Authenticator) logRefreshError(ctx context.Context, err error) {
	if aerr, ok := apperr.AsError(err); ok {
		statusCode := aerr.Code.HTTPStatusCode()
		if statusCode >= 400 && statusCode < 500 {
			slog.WarnContext(ctx, "Failed to refresh user token", "error", err)
			return
		}
	}
	slog.ErrorContext(ctx, "Failed to refresh user token", "error", err)
}
