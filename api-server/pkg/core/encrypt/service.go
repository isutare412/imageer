package encrypt

import (
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type service struct {
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

func NewService() Service {
	return &service{}
}
