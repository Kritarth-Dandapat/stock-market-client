// Package crypto provides AES-256-GCM encryption and HMAC-SHA256 signing
// utilities for protecting PII (PAN, Aadhaar, bank account numbers) before
// storing them in the database.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

var (
	// ErrInvalidKeyLength is returned when the encryption key is not 32 bytes.
	ErrInvalidKeyLength = errors.New("crypto: encryption key must be 32 bytes for AES-256")
	// ErrCiphertextTooShort is returned when the ciphertext is shorter than the GCM nonce.
	ErrCiphertextTooShort = errors.New("crypto: ciphertext too short")
)

// EncryptPII encrypts plaintext PII using AES-256-GCM.
// The returned string is base64-encoded and contains the random nonce
// prepended to the ciphertext.  The key must be exactly 32 bytes.
func EncryptPII(plaintext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKeyLength
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPII reverses EncryptPII, returning the original plaintext.
func DecryptPII(encoded string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKeyLength
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrCiphertextTooShort
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// SignRequest computes an HMAC-SHA256 signature for an outgoing BSE API
// request payload.  The result is hex-encoded.
func SignRequest(payload []byte, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature checks that the provided signature matches the payload
// using constant-time comparison to prevent timing attacks.
func VerifySignature(payload []byte, secret []byte, signature string) bool {
	expected := SignRequest(payload, secret)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// MaskPAN returns a masked version of a PAN card number, e.g. "ABCDE****F".
// If the PAN is shorter than 10 characters the original value is returned.
func MaskPAN(pan string) string {
	if len(pan) < 10 {
		return pan
	}
	return pan[:5] + "****" + pan[9:]
}

// MaskAadhaar returns a masked Aadhaar number showing only the last 4 digits.
func MaskAadhaar(aadhaar string) string {
	if len(aadhaar) < 4 {
		return "****"
	}
	masked := ""
	for i := 0; i < len(aadhaar)-4; i++ {
		masked += "*"
	}
	return masked + aadhaar[len(aadhaar)-4:]
}

// MaskBankAccount returns a masked bank account number showing only the last 4 digits.
func MaskBankAccount(account string) string {
	if len(account) < 4 {
		return "****"
	}
	masked := ""
	for i := 0; i < len(account)-4; i++ {
		masked += "*"
	}
	return masked + account[len(account)-4:]
}
