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
	nonceBuffer := make([]byte, 12)
	binary.BigEndian.PutUint32(nonceBuffer[:4], nonce)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Set the RTP header as additional authenticated data (AAD)
	ciphertext := aesgcm.Seal(nil, nonceBuffer, plaintext, rtpHeader)

	// Create a buffer to hold the final packet
	headerSize := len(rtpHeader)
	packetSize := headerSize + len(ciphertext) + 4 // 4 bytes for the nonce padding
	packet := make([]byte, packetSize)

	// Copy the RTP header to the packet
	copy(packet, rtpHeader)

	// Copy the encrypted data to the packet
	copy(packet[headerSize:], ciphertext)

	// Copy the first 4 bytes of the nonce to the end of the packet
	copy(packet[headerSize+len(ciphertext):], nonceBuffer[:4])

	return packet, nil
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
