package auth

import (
	"github.com/golang-jwt/jwt"
)

type Token string

type Session struct {
	Id        string
	Privilege string
}

type claims struct {
	jwt.StandardClaims
	Privilege string
}
