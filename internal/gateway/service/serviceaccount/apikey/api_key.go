package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/jxskiss/base62"

	"github.com/isutare412/imageer/pkg/apperr"
)

const (
	prefix      = "ak_"
	randCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randLen     = 48
	checksumLen = 16
)

type APIKey string

func New() APIKey {
	body := generateRandomString(randLen)
	checksum := calculateChecksum(body)
	return APIKey(prefix + body + checksum)
}

func ParseString(s string) (APIKey, error) {
	k := APIKey(s)
	if err := k.Validate(); err != nil {
		return "", fmt.Errorf("validating API key: %w", err)
	}
	return k, nil
}

func (k APIKey) String() string {
	return string(k)
}

func (k APIKey) Validate() error {
	// check length
	if len(k) != len(prefix)+randLen+checksumLen {
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("API key is invalid")
	}

	// check prefix
	if k[:len(prefix)] != prefix {
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("API key is invalid")
	}

	// check checksum
	body, checksum := k[len(prefix):len(k)-checksumLen], k[len(k)-checksumLen:]
	if calculateChecksum(string(body)) != string(checksum) {
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("API key is invalid")
	}
	return nil
}

func (k APIKey) Hash() string {
	hashBytes := sha256.Sum256([]byte(k))
	hash := base64.RawStdEncoding.EncodeToString(hashBytes[:])
	return hash
}

func generateRandomString(size int) string {
	result := make([]byte, size)
	charsetLen := big.NewInt(int64(len(randCharset)))

	for i := range size {
		rand.Text()
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			panic(fmt.Sprintf("unexpected rand.Int failure: %v", err))
		}
		result[i] = randCharset[num.Int64()]
	}

	return string(result)
}

func calculateChecksum(s string) string {
	hash := sha256.Sum256([]byte(s))
	checksum := base62.EncodeToString(hash[:])[:checksumLen]
	return checksum
}
