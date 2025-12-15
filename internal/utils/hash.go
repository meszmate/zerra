package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
)

func GetFileHash(f io.Reader) (string, *errx.Error) {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}

	return hex.EncodeToString(h.Sum(nil))[:16], nil
}
