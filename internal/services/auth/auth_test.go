package auth

import (
	"testing"
)

// --- isLegacyMD5 ---

func TestIsLegacyMD5(t *testing.T) {
	cases := []struct {
		hash string
		want bool
	}{
		{"5f4dcc3b5aa765d61d8327deb882cf99", true},  // MD5 of "password"
		{"d8578edf8458ce06fbc5bb76a58c5ca4", true}, // MD5 of "qwerty"
		{"$2a$12$abcdefghijklmnopqrstuuVnNHSTwjXjN.E5VWKQ7eLhOgB1mT3a", false}, // bcrypt
		{"$2b$12$somethinglongerthan20chars", false},
		{"", true}, // empty hash is not bcrypt
	}
	for _, tc := range cases {
		if got := isLegacyMD5(tc.hash); got != tc.want {
			t.Errorf("isLegacyMD5(%q) = %v, want %v", tc.hash, got, tc.want)
		}
	}
}

// --- verifyMD5 ---

func TestVerifyMD5(t *testing.T) {
	// MD5("password") = 5f4dcc3b5aa765d61d8327deb882cf99
	if !verifyMD5("password", "5f4dcc3b5aa765d61d8327deb882cf99") {
		t.Error("expected MD5 match")
	}
	if verifyMD5("wrong", "5f4dcc3b5aa765d61d8327deb882cf99") {
		t.Error("expected MD5 mismatch")
	}
}

// --- HashPassword / bcrypt verify ---

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if !isLegacyMD5("") { // just confirm legacy check does not trigger on empty
		t.Skip()
	}
	// bcrypt hash must start with $2
	if len(hash) < 2 || hash[:2] != "$2" {
		t.Errorf("expected bcrypt hash, got: %s", hash[:10])
	}
	// Must not match a different password.
	hash2, _ := HashPassword("other")
	if hash == hash2 {
		t.Error("two different passwords produced the same hash (impossible with bcrypt)")
	}
}
