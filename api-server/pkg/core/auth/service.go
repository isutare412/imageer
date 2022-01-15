package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	"github.com/isutare412/imageer/api-server/pkg/config"
)

type Service interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool

	SignToken(id ID) (Token, error)
	VerifyToken(t Token) (ID, error)
}

type service struct {
	signKey    *rsa.PrivateKey
	verifyKey  *rsa.PublicKey
	expireHour int64
}

func (s *service) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *service) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *service) SignToken(id ID) (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":  string(id),
		"exp": time.Now().Add(time.Duration(s.expireHour) * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString(s.signKey)
	if err != nil {
		return "", err
	}
	return Token(tokenStr), nil
}

func (s *service) VerifyToken(t Token) (ID, error) {
	token, err := jwt.Parse(string(t), func(t *jwt.Token) (interface{}, error) {
		return s.verifyKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("token is not jwt.MapClaims")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", ErrTokenExpired
	}

	idInterface, ok := claims["id"]
	if !ok {
		return "", errors.New("id not found in claims")
	}
	id, ok := idInterface.(string)
	if !ok {
		return "", errors.New("id claim is not string")
	}
	return ID(id), nil
}

func ContextWithID(ctx context.Context, id ID) context.Context {
	return context.WithValue(ctx, ctxKeyID, id)
}

func IDFromContext(ctx context.Context) (ID, error) {
	val := ctx.Value(ctxKeyID)
	if val == nil {
		return "", ErrCtxIDNotFound
	}
	id, ok := val.(ID)
	if !ok {
		return "", ErrCtxInvalidID
	}
	return id, nil
}

func NewService(cfg *config.AuthConfig) (Service, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
	if err != nil {
		return nil, err
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
	if err != nil {
		return nil, err
	}

	if cfg.ExpireHour <= 0 {
		return nil, fmt.Errorf("ExpireHour[%d] should be greater than 0", cfg.ExpireHour)
	}

	return &service{
		expireHour: cfg.ExpireHour,
		signKey:    signKey,
		verifyKey:  verifyKey,
	}, nil
}
