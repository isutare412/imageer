package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/isutare412/imageer/pkg/apperr"
)

// AESCrypter implements encryption/decryption using AES-256-GCM encryption.
type AESCrypter struct {
	gcm cipher.AEAD
}

// NewAESCrypter creates a new AESCrypter with the provided 32-byte key.
func NewAESCrypter(cfg AESCrypterConfig) (*AESCrypter, error) {
	key := cfg.KeyBytes()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to create AES cipher").
			WithCause(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to create GCM cipher").
			WithCause(err)
	}

	return &AESCrypter{gcm: gcm}, nil
}

// Encrypt encrypts the given data using AES-256-GCM.
// The returned data contains both the nonce and the encrypted data.
func (c *AESCrypter) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to generate nonce").
			WithCause(err)
	}

	ciphertext := c.gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts the given encrypted data using AES-256-GCM.
// The encrypted data must contain both the nonce and the encrypted data.
func (c *AESCrypter) Decrypt(encryptedData []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("encrypted data too short: expected at least %d bytes, got %d", nonceSize, len(encryptedData))
	}

	nonce := encryptedData[:nonceSize]
	ciphertext := encryptedData[nonceSize:]

	plaintext, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("failed to decrypt data").
			WithCause(err)
	}

	return plaintext, nil
}
