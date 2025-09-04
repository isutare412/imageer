package crypt

import "crypto/sha256"

type AESCrypterConfig struct {
	Key string
}

// KeyBytes returns the 32-byte key derived from key string.
func (c AESCrypterConfig) KeyBytes() []byte {
	hash := sha256.Sum256([]byte(c.Key))
	return hash[:]
}
