package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Generate256BToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
