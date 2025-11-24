package auth

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/meszmate/zerra/internal/config"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/notify/templates"
	"github.com/meszmate/zerra/internal/pkg/crypt"
)

func (s *authService) ResetPasswordStart(ctx context.Context, data *ResetPasswordStart, ipaddr string) *errx.Error {
	if err := s.captcha.Verify(ctx, data.Turnstile, ipaddr); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	user, err := s.userRepostory.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return err
	}

	u, err := s.userService.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}

	if err := s.passwordResetLimit(ctx, u.Email); err != nil {
		return err
	}

	sessionID := uuid.NewString()
	nonce, xerr := crypt.Nonce()
	if xerr != nil {
		sentry.CaptureException(xerr)
		return errx.InternalError()
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(PasswordResetTTL)

	token, xerr := s.tokenService.GenerateToken(user.ID, sessionID, nonce, issuedAt, expiresAt)
	if xerr != nil {
		sentry.CaptureException(xerr)
		return errx.InternalError()
	}

	if err := s.saveResetPasswordSession(ctx, sessionID, nonce); err != nil {
		return err
	}

	url := config.GetPasswordResetURL(token)

	text, xerr := templates.GenerateResetPasswordHTML(u.FirstName, url)
	if xerr != nil {
		return errx.InternalError()
	}

	if err := s.emailNotificationService.Send(ctx, []string{u.Email}, nil, nil, "Password Reset Confirmation", text); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *authService) ResetPasswordConfirm(ctx context.Context, data *ResetPasswordConfirm, session, ipaddr string) *errx.Error {
	if err := s.captcha.Verify(ctx, data.Turnstile, ipaddr); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	sess, err := s.tokenService.VerifyToken(session)
	if err != nil {
		return err
	}

	if sess.ExpiresAt.Before(time.Now()) {
		return errx.ErrToken
	}

	nonce, err := s.getResetPasswordSession(ctx, sess.SessionID)
	if err != nil {
		return err
	}

	if nonce != sess.Nonce {
		return errx.ErrToken
	}

	if err := s.deletePasswordResetSession(ctx, sess.SessionID); err != nil {
		return err
	}

	if err := s.authRepostory.ResetPassword(ctx, sess.UserID, data.Password); err != nil {
		return err
	}

	return nil
}
