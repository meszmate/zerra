package crypt

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

func SHA256(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func HashToRange(input string, max int) int {
	h := sha256.Sum256([]byte(input))
	num := binary.BigEndian.Uint64(h[:8])
	return int(num % uint64(max+1))
}
