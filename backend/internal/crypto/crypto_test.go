package crypto_test

import (
	"strings"
	"testing"

	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/crypto"
)

func TestEncryptDecryptPII(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 1)
	}

	plaintext := "ABCDE1234F"
	encrypted, err := crypto.EncryptPII(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptPII error: %v", err)
	}
	if encrypted == plaintext {
		t.Fatal("encrypted text must not equal plaintext")
	}

	decrypted, err := crypto.DecryptPII(encrypted, key)
	if err != nil {
		t.Fatalf("DecryptPII error: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptPIIInvalidKeyLength(t *testing.T) {
	_, err := crypto.EncryptPII("test", []byte("shortkey"))
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestDecryptPIIDifferentKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	for i := range key1 {
		key1[i] = byte(i + 1)
		key2[i] = byte(i + 2)
	}

	encrypted, err := crypto.EncryptPII("secret", key1)
	if err != nil {
		t.Fatalf("EncryptPII error: %v", err)
	}

	_, err = crypto.DecryptPII(encrypted, key2)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong key")
	}
}

func TestEncryptPIINonDeterministic(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 1)
	}

	e1, _ := crypto.EncryptPII("test", key)
	e2, _ := crypto.EncryptPII("test", key)
	if e1 == e2 {
		t.Fatal("each encryption should produce a unique ciphertext due to random nonce")
	}
}

func TestSignAndVerifyRequest(t *testing.T) {
	secret := []byte("my-hmac-secret-key")
	payload := []byte(`{"amount":5000,"scheme":"INF123456789"}`)

	sig := crypto.SignRequest(payload, secret)
	if sig == "" {
		t.Fatal("expected non-empty signature")
	}

	if !crypto.VerifySignature(payload, secret, sig) {
		t.Fatal("signature verification should succeed for matching payload")
	}

	if crypto.VerifySignature([]byte("tampered"), secret, sig) {
		t.Fatal("signature verification must fail for tampered payload")
	}
}

func TestMaskPAN(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"ABCDE1234F", "ABCDE****F"},
		{"SHORT", "SHORT"},
		{"", ""},
	}
	for _, c := range cases {
		got := crypto.MaskPAN(c.input)
		if got != c.expected {
			t.Errorf("MaskPAN(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestMaskAadhaar(t *testing.T) {
	got := crypto.MaskAadhaar("123456789012")
	if !strings.HasSuffix(got, "9012") {
		t.Errorf("MaskAadhaar should end with last 4 digits, got %q", got)
	}
	if strings.Contains(got[:len(got)-4], "1") {
		t.Errorf("MaskAadhaar should mask all but last 4 digits, got %q", got)
	}
}

func TestMaskBankAccount(t *testing.T) {
	got := crypto.MaskBankAccount("123456789")
	if !strings.HasSuffix(got, "6789") {
		t.Errorf("MaskBankAccount should end with last 4 digits, got %q", got)
	}
}
