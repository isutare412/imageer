package immigration

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

type PassportIssuer struct {
	authSvc           port.AuthService
	serviceAccountSvc port.ServiceAccountService
	apiKeyHeader      string
	userCookieName    string
}

func NewPassportIssuer(apiKeyHeader, userCookieName string, authSvc port.AuthService,
	serviceAccountSvc port.ServiceAccountService,
) *PassportIssuer {
	return &PassportIssuer{
		authSvc:           authSvc,
		serviceAccountSvc: serviceAccountSvc,
		apiKeyHeader:      apiKeyHeader,
		userCookieName:    userCookieName,
	}
}

func (i *PassportIssuer) IssuePassport(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rctx := ctx.Request().Context()

		if ok, err := i.issuePassportByHeader(rctx, ctx.Request().Header); err != nil {
			return fmt.Errorf("issue passport by header: %w", err)
		} else if ok {
			// Passport issued by header
			return next(ctx)
		}

		if ok, err := i.issuePassportByCookie(rctx, ctx.Cookies()); err != nil {
			return fmt.Errorf("issue passport by cookie: %w", err)
		} else if ok {
			// Passport issued by cookie
			return next(ctx)
		}

		// No passport issued
		return next(ctx)
	}
}

func (i *PassportIssuer) issuePassportByHeader(ctx context.Context, header http.Header) (bool, error) {
	if ok, err := i.passportFromBearerToken(ctx, header); err != nil {
		return false, fmt.Errorf("getting passport from bearer token: %w", err)
	} else if ok {
		return true, nil
	}

	if ok, err := i.passportFromAPIKey(ctx, header); err != nil {
		return false, fmt.Errorf("getting passport from API key: %w", err)
	} else if ok {
		return true, nil
	}

	return false, nil
}

func (i *PassportIssuer) passportFromBearerToken(ctx context.Context, header http.Header) (bool, error) {
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

	passport := domain.NewUserTokenPassport(payload)
	i.registerPassport(ctx, passport)
	return true, nil
}

func (i *PassportIssuer) passportFromAPIKey(ctx context.Context, header http.Header) (bool, error) {
	apiKey := header.Get(i.apiKeyHeader)
	if apiKey == "" {
		return false, nil
	}

	account, err := i.serviceAccountSvc.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return false, fmt.Errorf("getting service account by API key: %w", err)
	}

	passport := domain.NewServiceAccountPassport(account)
	i.registerPassport(ctx, passport)
	return true, nil
}

func (i *PassportIssuer) issuePassportByCookie(ctx context.Context, cookies []*http.Cookie) (bool, error) {
	cookie, ok := lo.Find(cookies, func(c *http.Cookie) bool { return c.Name == i.userCookieName })
	if !ok {
		return false, nil
	}

	payload, err := i.authSvc.VerifyUserToken(ctx, cookie.Value)
	if err != nil {
		return false, fmt.Errorf("verifying user token: %w", err)
	}

	passport := domain.NewUserTokenPassport(payload)
	i.registerPassport(ctx, passport)
	return true, nil
}

func (i *PassportIssuer) registerPassport(ctx context.Context, pp domain.Passport) {
	if bag, ok := contextbag.BagFromContext(ctx); ok {
		bag.Passport = pp
	}
}
