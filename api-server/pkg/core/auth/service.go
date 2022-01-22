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

	SignToken(sess *Session) (Token, error)
	VerifyToken(t Token) (*Session, error)
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

func (s *service) SignToken(sess *Session) (Token, error) {
	now := time.Now()
	expire := now.Add(time.Duration(s.expireHour) * time.Hour)
	clm := claims{
		StandardClaims: jwt.StandardClaims{
			Id:        sess.Id,
			IssuedAt:  now.Unix(),
			ExpiresAt: expire.Unix(),
		},
		Privilege: sess.Privilege,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &clm)
	tokenStr, err := token.SignedString(s.signKey)
	if err != nil {
		return "", err
	}
	return Token(tokenStr), nil
}

func (s *service) VerifyToken(t Token) (*Session, error) {
	token, err := jwt.ParseWithClaims(string(t), &claims{}, func(t *jwt.Token) (interface{}, error) {
		return s.verifyKey, nil
	})
	if err != nil {
		return nil, err
	}

	clm, ok := token.Claims.(*claims)
	if !ok {
		return nil, errors.New("token is not claims")
	}
	if !clm.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, ErrTokenExpired
	}

	sess := Session{
		Id:        clm.Id,
		Privilege: clm.Privilege,
	}
	return &sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, ctxKeyID, sess)
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	val := ctx.Value(ctxKeyID)
	if val == nil {
		return nil, ErrCtxSessionNotFound
	}
	sess, ok := val.(*Session)
	if !ok {
		return nil, ErrCtxInvalidSession
	}
	return sess, nil
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
