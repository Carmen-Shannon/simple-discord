package crypto

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

// EncryptXChaCha20Poly1305 encrypts the plaintext using XChaCha20-Poly1305
func EncryptXChaCha20Poly1305(plaintext, key, rtpHeader []byte, nonce uint32) ([]byte, error) {
	// Convert the 32-bit incremental nonce value into a 24-byte slice
	nonceBuffer := make([]byte, 24)
	binary.BigEndian.PutUint32(nonceBuffer[:4], nonce)

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create XChaCha20-Poly1305: %w", err)
	}

	ciphertext := aead.Seal(nil, nonceBuffer, plaintext, rtpHeader)

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
