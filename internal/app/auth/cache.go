package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/crypt"
	"github.com/redis/go-redis/v9"
)

func getEmailVerificationKey(email string) string {
	return "email_verification:" + crypt.SHA256(email)
}

func getPasswordResetLimitKey(email string) string {
	return "password_reset_limit:" + crypt.SHA256(email)
}

func getLoginSessionKey(sessionID string) string {
	return "login_sess:" + sessionID
}

func getRegistrationSessionKey(sessionID string) string {
	return "registration_sess:" + sessionID
}

func getResetPasswordSessionKey(sessionID string) string {
	return "reset_password:" + sessionID
}

func (s *authService) saveLoginSession(ctx context.Context, sessionID string, session *models.LoginSession, expiresAt time.Time) *errx.Error {
	data, err := json.Marshal(session)
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if err := s.cache.Set(ctx, getLoginSessionKey(sessionID), data, time.Until(expiresAt)).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *authService) getLoginSession(ctx context.Context, sessionID string) (*models.LoginSession, *errx.Error) {
	data, err := s.cache.Get(ctx, getLoginSessionKey(sessionID)).Bytes()
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

func (s *authService) saveRegistrationSession(ctx context.Context, sessionID string, session *models.RegistrationSession, expiresAt time.Time) *errx.Error {
	data, err := json.Marshal(session)
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if err := s.cache.Set(ctx, getRegistrationSessionKey(sessionID), data, time.Until(expiresAt)).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *authService) getRegistrationSession(ctx context.Context, sessionID string) (*models.RegistrationSession, *errx.Error) {
	data, err := s.cache.Get(ctx, getRegistrationSessionKey(sessionID)).Bytes()
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

	count, err := s.cache.Incr(ctx, key).Result()
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if count == 1 {
		if err := s.cache.Expire(ctx, key, AuthEmailTTL).Err(); err != nil {
			sentry.CaptureException(err)
			return errx.InternalError()
		}
	}

	if count > AuthEmailLimit {
		return errx.ErrAuthLimit
	}

	return nil
}

func (s *authService) passwordResetLimit(ctx context.Context, email string) *errx.Error {
	key := getPasswordResetLimitKey(email)

	count, err := s.cache.Incr(ctx, key).Result()
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if count == 1 {
		if err := s.cache.Expire(ctx, key, PasswordResetLimitTTL).Err(); err != nil {
			sentry.CaptureException(err)
			return errx.InternalError()
		}
	}

	if count > PasswordResetLimit {
		return errx.ErrAuthLimit
	}

	return nil
}

func (s *authService) saveResetPasswordSession(ctx context.Context, sessionID string, nonce string) *errx.Error {
	if err := s.cache.SetEx(ctx, getResetPasswordSessionKey(sessionID), nonce, SessionTTL).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *authService) getResetPasswordSession(ctx context.Context, sessionID string) (string, *errx.Error) {
	val, err := s.cache.Get(ctx, getResetPasswordSessionKey(sessionID)).Result()
	if err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}

	return val, nil
}

func (s *authService) deletePasswordResetSession(ctx context.Context, sessionID string) *errx.Error {
	val, err := s.cache.Del(ctx, getResetPasswordSessionKey(sessionID)).Result()
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if val == 0 {
		return errx.ErrToken
	}

	return nil
}
