package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/users"
)

const imageerGatewayIssuer = "imageer-gateway"

type appClaims struct {
	jwt.RegisteredClaims
	UserID     string     `json:"user_id"`
	Role       users.Role `json:"role"`
	Nickname   string     `json:"nickname"`
	Email      string     `json:"email"`
	PictureURL string     `json:"picture_url"`
}

func newAppClaims(payload domain.UserTokenPayload) appClaims {
	return appClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    imageerGatewayIssuer,
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			NotBefore: jwt.NewNumericDate(payload.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(payload.ExpireAt),
		},
		UserID:     payload.UserID,
		Role:       payload.Role,
		Nickname:   payload.Nickname,
		Email:      payload.Email,
		PictureURL: payload.PictureURL,
	}
}

func (c *appClaims) toUserTokenPayload() domain.UserTokenPayload {
	return domain.UserTokenPayload{
		UserID:     c.UserID,
		IssuedAt:   c.IssuedAt.Time,
		ExpireAt:   c.ExpiresAt.Time,
		Role:       c.Role,
		Nickname:   c.Nickname,
		Email:      c.Email,
		PictureURL: c.PictureURL,
	}
}
