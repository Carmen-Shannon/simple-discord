package crypto

import (
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

// EncryptXChaCha20Poly1305 encrypts the plaintext using XChaCha20-Poly1305 and uses the RTP header as the nonce
func EncryptXChaCha20Poly1305(plaintext, key, nonce []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create XChaCha20-Poly1305: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptXChaCha20Poly1305 decrypts the ciphertext using XChaCha20-Poly1305 and uses the RTP header as the nonce
func DecryptXChaCha20Poly1305(ciphertext, key, nonce []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create XChaCha20-Poly1305: %w", err)
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt XChaCha20-Poly1305: %w", err)
	}

	return plaintext, nil
}
