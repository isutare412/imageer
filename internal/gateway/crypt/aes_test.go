package crypt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAESCrypter(t *testing.T) {
	tests := []struct {
		name      string
		cfg       AESCrypterConfig
		expectErr bool
	}{
		{
			name:      "valid key string",
			cfg:       AESCrypterConfig{Key: "test-key"},
			expectErr: false,
		},
		{
			name:      "empty key string",
			cfg:       AESCrypterConfig{Key: ""},
			expectErr: false,
		},
		{
			name:      "long key string",
			cfg:       AESCrypterConfig{Key: "this-is-a-very-long-key-string-for-testing"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crypter, err := NewAESCrypter(tt.cfg)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, crypter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, crypter)
			}
		})
	}
}

func TestAESCrypter_EncryptDecrypt(t *testing.T) {
	cfg := AESCrypterConfig{Key: "test-key-for-encryption"}
	crypter, err := NewAESCrypter(cfg)
	require.NoError(t, err)

	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "small text",
			data: []byte("hello world"),
		},
		{
			name: "longer text",
			data: []byte("This is a longer text that should be encrypted and decrypted properly."),
		},
		{
			name: "binary data",
			data: []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
		},
		{
			name: "unicode text",
			data: []byte("Hello, 世界! 🌍"),
		},
		{
			name: "large data",
			data: []byte(strings.Repeat("abcdefghijklmnopqrstuvwxyz", 1000)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the data
			encrypted, err := crypter.Encrypt(tt.data)
			require.NoError(t, err)
			assert.NotEqual(t, tt.data, encrypted)
			assert.Greater(t, len(encrypted), len(tt.data))

			// Decrypt the data
			decrypted, err := crypter.Decrypt(encrypted)
			require.NoError(t, err)
			if len(tt.data) == 0 {
				assert.Empty(t, decrypted)
			} else {
				assert.Equal(t, tt.data, decrypted)
			}
		})
	}
}

func TestAESCrypter_DecryptErrors(t *testing.T) {
	cfg := AESCrypterConfig{Key: "test-key-for-errors"}
	crypter, err := NewAESCrypter(cfg)
	require.NoError(t, err)

	tests := []struct {
		name          string
		encryptedData []byte
		expectErr     bool
	}{
		{
			name:          "empty data",
			encryptedData: []byte{},
			expectErr:     true,
		},
		{
			name:          "too short data",
			encryptedData: []byte{0x01, 0x02},
			expectErr:     true,
		},
		{
			name:          "invalid ciphertext",
			encryptedData: make([]byte, 20), // 12 bytes nonce + 8 bytes invalid ciphertext
			expectErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := crypter.Decrypt(tt.encryptedData)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAESCrypter_DifferentKeysProduceDifferentResults(t *testing.T) {
	cfg1 := AESCrypterConfig{Key: "test-key-1"}
	cfg2 := AESCrypterConfig{Key: "test-key-2"}

	crypter1, err := NewAESCrypter(cfg1)
	require.NoError(t, err)
	crypter2, err := NewAESCrypter(cfg2)
	require.NoError(t, err)

	plaintext := []byte("test data")

	encrypted1, err := crypter1.Encrypt(plaintext)
	require.NoError(t, err)
	encrypted2, err := crypter2.Encrypt(plaintext)
	require.NoError(t, err)

	// Different keys should produce different encrypted results
	assert.NotEqual(t, encrypted1, encrypted2)

	// Each crypter should only be able to decrypt its own encrypted data
	decrypted1, err := crypter1.Decrypt(encrypted1)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted1)

	_, err = crypter2.Decrypt(encrypted1)
	assert.Error(t, err, "crypter2 should not be able to decrypt data encrypted with crypter1")
}

func TestAESCrypter_SameDataProducesDifferentCiphertext(t *testing.T) {
	cfg := AESCrypterConfig{Key: "test-key-for-same-data"}
	crypter, err := NewAESCrypter(cfg)
	require.NoError(t, err)

	plaintext := []byte("test data")

	// Encrypt the same data multiple times
	encrypted1, err := crypter.Encrypt(plaintext)
	require.NoError(t, err)
	encrypted2, err := crypter.Encrypt(plaintext)
	require.NoError(t, err)
	encrypted3, err := crypter.Encrypt(plaintext)
	require.NoError(t, err)

	// Each encryption should produce different ciphertext due to random nonces
	assert.NotEqual(t, encrypted1, encrypted2)
	assert.NotEqual(t, encrypted2, encrypted3)
	assert.NotEqual(t, encrypted1, encrypted3)

	// But all should decrypt to the same plaintext
	decrypted1, err := crypter.Decrypt(encrypted1)
	require.NoError(t, err)
	decrypted2, err := crypter.Decrypt(encrypted2)
	require.NoError(t, err)
	decrypted3, err := crypter.Decrypt(encrypted3)
	require.NoError(t, err)

	assert.Equal(t, plaintext, decrypted1)
	assert.Equal(t, plaintext, decrypted2)
	assert.Equal(t, plaintext, decrypted3)
}

func BenchmarkAESCrypter_Encrypt(b *testing.B) {
	cfg := AESCrypterConfig{Key: "benchmark-key"}
	crypter, err := NewAESCrypter(cfg)
	require.NoError(b, err)

	data := []byte("This is some test data to encrypt for benchmarking purposes.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypter.Encrypt(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAESCrypter_Decrypt(b *testing.B) {
	cfg := AESCrypterConfig{Key: "benchmark-key"}
	crypter, err := NewAESCrypter(cfg)
	require.NoError(b, err)

	data := []byte("This is some test data to decrypt for benchmarking purposes.")
	encrypted, err := crypter.Encrypt(data)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypter.Decrypt(encrypted)
		if err != nil {
			b.Fatal(err)
		}
	}
}
