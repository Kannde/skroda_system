package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// GenerateReferenceCode creates a human-readable transaction reference.
// Format: SKR-2026-A7X9 (no 0/O/1/I to avoid confusion)
func GenerateReferenceCode() string {
	year := time.Now().Year()
	chars := "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 4)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[n.Int64()]
	}
	return fmt.Sprintf("SKR-%d-%s", year, string(code))
}

// GenerateInviteToken creates a secure random token for transaction invites.
func GenerateInviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
