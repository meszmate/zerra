package crypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
