package auth

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/pkg/crypt"
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

func (s *authService) EmailVerificationLimit(ctx context.Context, email string) (int64, *errx.Error) {
	key := getEmailVerificationKey(email)
	cmd := s.Cache.Incr(ctx, key)
	if err := cmd.Err(); err != nil {
		sentry.CaptureException(err)
		return 0, errx.InternalError()
	}

	val := cmd.Val()
	if val == 1 {
		if err := s.Cache.Expire(ctx, key, AuthEmailTTL).Err(); err != nil {
			sentry.CaptureException(err)
			return 0, errx.InternalError()
		}
	}

	if val > AuthEmailLimit {
		return 0, errx.ErrEmailLimit
	}

	return val, nil
}

func (s *authService) GenerateLoginSession(ctx context.Context, sessionID, userID, nonce string) *errx.Error{
	key := getLoginSessionKey(sessionID)

	if err := s.Cache.SetEx(ctx, key, nonce, AuthSessionTTL).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *authService) GetLoginSession(ctx context.Context, sessionID string) (string, *errx.Error) {
	key := getLoginSessionKey(sessionID)

	cmd := s.Cache.Get(ctx, key);
	if err := cmd.Err(); err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}

	return cmd.Val(), nil
}


func (s *authService) GenerateRegistrationSession(ctx context.Context, sessionID, userID, nonce string) *errx.Error{
	key := getLoginSessionKey(sessionID)

	if err := s.Cache.SetEx(ctx, key, nonce, AuthSessionTTL).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *authService) GetLoginSession(ctx context.Context, sessionID) (string, *errx.Error) {
	key := getLoginSessionKey(sessionID)

	cmd := s.Cache.Get(ctx, key);
	if err := cmd.Err(); err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}

	return cmd.Val(), nil
}
