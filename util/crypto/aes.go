package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
)

// EncryptAESGCM encrypts the plaintext using AES-GCM and uses the RTP header as additional data
func EncryptAESGCM(plaintext, key []byte, rtpHeader []byte, nonce uint32) ([]byte, error) {
	// Convert the 32-bit incremental nonce value into a 12-byte slice
	nonceBytes := make([]byte, 12)
	binary.LittleEndian.PutUint32(nonceBytes[:4], nonce)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Set the RTP header as additional authenticated data (AAD)
	ciphertext := aesgcm.Seal(nil, nonceBytes, plaintext, rtpHeader)

	// Append the 32-bit incremental nonce value to the encrypted payload
	ciphertext = append(ciphertext, nonceBytes[:4]...)
	return ciphertext, nil
}

// DecryptAESGCM decrypts the ciphertext using AES-GCM and uses the RTP header as the nonce
func DecryptAESGCM(ciphertext, key, nonce, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt AES GCM: %w", err)
	}

	return plaintext, nil
}
