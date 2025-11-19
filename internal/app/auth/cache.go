package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/crypt"
	"github.com/redis/go-redis/v9"
)

func getEmailVerificationKey(email string) string {
	return "email_verification:" + crypt.SHA256(email)
}

func getLoginSessionKey(sessionID string) string {
	return "login_sess:" + sessionID
}

func getRegistrationSessionKey(sessionID string) string {
	return "registration_sess:" + sessionID
}

func (s *authService) GenerateLoginSession(ctx context.Context, sessionID string, data *models.LoginSession) *errx.Error {
	key := getLoginSessionKey(sessionID)

	if err := s.Cache.SetEx(ctx, key, data, AuthSessionTTL).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *authService) GetLoginSession(ctx context.Context, sessionID string) (*models.LoginSession, *errx.Error) {
	data, err := s.Cache.Get(ctx, getLoginSessionKey(sessionID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	var session models.LoginSession
	if err := json.Unmarshal(data, &session); err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	return &session, nil
}

func (s *authService) GenerateRegistrationSession(ctx context.Context, sessionID, userID, nonce string) *errx.Error {
	key := getLoginSessionKey(sessionID)

	if err := s.Cache.SetEx(ctx, key, nonce, AuthSessionTTL).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *authService) GetRegistrationSession(ctx context.Context, sessionID string) (*models.RegistrationSession, *errx.Error) {
	data, err := s.Cache.Get(ctx, getRegistrationSessionKey(sessionID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	var session models.RegistrationSession
	if err := json.Unmarshal(data, &session); err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	return &session, nil
}

func (s *authService) canSendEmail(ctx context.Context, email string) *errx.Error {
	key := getEmailVerificationKey(email)

	count, err := s.Cache.Incr(ctx, key).Result()
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if count == 1 {
		if err := s.Cache.Expire(ctx, key, AuthEmailTTL).Err(); err != nil {
			sentry.CaptureException(err)
			return errx.InternalError()
		}
	}

	if count > AuthEmailLimit {
		return errx.ErrEmailLimit
	}

	return nil
}
