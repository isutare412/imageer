package oidc

import (
	"context"
	"fmt"
	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/oauth2"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/tracing"
)

type GoogleClient struct {
	oidcProvider *oidc.Provider
	oidcVerifier *oidc.IDTokenVerifier
	oauthCfg     *oauth2.Config
	cfg          GoogleClientConfig
}

func NewGoogleClient(cfg GoogleClientConfig) (*GoogleClient, error) {
	oidcProvider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("creating OIDC provider: %w", err)
	}

	oidcVerifier := oidcProvider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return &GoogleClient{
		oidcProvider: oidcProvider,
		oidcVerifier: oidcVerifier,
		oauthCfg: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     oidcProvider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		},
		cfg: cfg,
	}, nil
}

func (c *GoogleClient) BuildAuthenticationURL(baseURL, state string) string {
	return c.oauthCfg.AuthCodeURL(state,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("nonce", uuid.NewString()),
		oauth2.SetAuthURLParam("prompt", "consent select_account"),
		oauth2.SetAuthURLParam("redirect_uri", c.redirectURI(baseURL)))
}

func (c *GoogleClient) ExchangeCode(ctx context.Context, baseURL, code string,
) (payload domain.IDTokenPayload, err error) {
	ctx, span := tracing.StartSpan(ctx, "oidc.GoogleClient.ExchangeCode")
	defer span.End()

	token, err := c.oauthCfg.Exchange(ctx, code,
		oauth2.SetAuthURLParam("redirect_uri", c.redirectURI(baseURL)))
	if err != nil {
		return payload, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("failed to exchange code for token")
	}

	idTokenRaw, ok := token.Extra("id_token").(string)
	if !ok {
		return payload, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("token response has invalid id_token")
	}

	idToken, err := c.oidcVerifier.Verify(ctx, idTokenRaw)
	if err != nil {
		return payload, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("failed to verify ID token")
	}

	idTokenPayload, err := googleIDTokenToDomain(idToken)
	if err != nil {
		return payload, fmt.Errorf("building Google ID token: %w", err)
	}

	return idTokenPayload, nil
}

func (c *GoogleClient) redirectURI(baseURL string) string {
	return lo.Must(url.JoinPath(baseURL, c.cfg.CallbackPath))
}
