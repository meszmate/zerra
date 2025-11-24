package auth

import (
	"context"
	"net/mail"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/argon2"
	"github.com/meszmate/zerra/internal/pkg/crypt"
	"github.com/meszmate/zerra/internal/pkg/gen"
)

func (s *authService) RegistrationStart(ctx context.Context, data *AuthData, ipaddr string) (*models.AuthSession, *errx.Error) {
	if xerr := s.captcha.Verify(ctx, data.Turnstile, ipaddr); xerr != nil {
		sentry.CaptureException(xerr)
		return nil, errx.InternalError()
	}

	if err := s.canSendEmail(ctx, data.Email); err != nil {
		return nil, err
	}

	passwordHash, xerr := argon2.Hash(data.Password)
	if xerr != nil {
		sentry.CaptureException(xerr)
		return nil, errx.InternalError()
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(AuthSessionTTL)
	sessionID := uuid.NewString()
	nonce, xerr := crypt.Nonce()
	if xerr != nil {
		sentry.CaptureException(xerr)
		return nil, errx.InternalError()
	}

	code, xerr := gen.VerificationCode()
	if xerr != nil {
		sentry.CaptureException(xerr)
		return nil, errx.InternalError()
	}

	codeHash, xerr := argon2.Hash(code)
	if xerr != nil {
		sentry.CaptureException(xerr)
		return nil, errx.InternalError()
	}

	session := &models.RegistrationSession{
		CodeHash:     codeHash,
		PasswordHash: passwordHash,
		Nonce:        nonce,
	}

	if err := s.saveRegistrationSession(ctx, sessionID, session, expiresAt); err != nil {
		return nil, err
	}

	sessionToken, xerr := s.tokenService.GenerateToken(data.Email, sessionID, nonce, issuedAt, expiresAt)
	if xerr != nil {
		sentry.CaptureException(xerr)
		return nil, errx.InternalError()
	}

	return &models.AuthSession{
		Session: sessionToken,
	}, nil
}

func (s *authService) RegistrationConfirm(ctx context.Context, data *ConfirmData, session, ipaddr string) *errx.Error {
	token, err := s.tokenService.VerifyToken(session)
	if err != nil {
		return err
	}
	if token.ExpiresAt.Before(time.Now()) {
		return errx.ErrSession
	}
	sess, err := s.getRegistrationSession(ctx, token.SessionID)
	if err != nil {
		return err
	}
	if sess == nil || sess.Nonce != token.Nonce {
		return errx.ErrSession
	}

	if sess.Tries >= AuthAttempts {
		return errx.ErrCodeLimit
	}

	v, xerr := argon2.Verify(data.Code, sess.CodeHash)
	if xerr != nil {
		sentry.CaptureException(xerr)
		return errx.InternalError()
	}

	if !v {
		return errx.ErrCode
	}

	email, xerr := mail.ParseAddress(token.UserID)
	if xerr != nil {
		return errx.ErrEmail
	}

	u, err := s.userRepostory.CreateUser(ctx, token.UserID, email.Name)
	if err != nil {
		return err
	}

	if err := s.userService.SaveUser(ctx, u); err != nil {
		return err
	}

	return nil
}
