package auth

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/argon2"
)

func (s *authService) GenerateLoginSession(ctx context.Context, userID, code string) (string, *errx.Error) {
	sessID := uuid.NewString()

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(SessionTTL)

	codeHash, err := argon2.Hash(code)
	if err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}

	loginSession := &models.LoginSession{
		CodeHash: codeHash,
		Tries:    0,
	}

	if err := s.saveLoginSession(ctx, sessID, loginSession, expiresAt); err != nil {
		return "", err
	}

	sessionToken, err := s.tokenService.GenerateToken(userID, sessID, "", issuedAt, expiresAt)
	if err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}

	return sessionToken, nil
}
