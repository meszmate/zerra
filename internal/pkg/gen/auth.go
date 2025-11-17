package gen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

const (
	AUTHORIZATION_TOKEN_LEN = 64
	AUTH_SESSION_TOKEN_LEN  = AUTHORIZATION_TOKEN_LEN
)

func VerificationCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func Token(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func IsValidHex(s string, length int) bool {
	if len(s) != length*2 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
