package token

import (
	"context"
	"time"

	"github.com/meszmate/zerra/internal/errx"
)

func (s *tokenService) GetSession(ctx context.Context, sessionID string) (*models.Session, *errx.Error) {
	sess, err := s.getSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if sess != nil {
		return sess, nil
	}

	sess, err = s.tokenRepostory.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.saveSession(ctx, sess, SessionTTL); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *tokenService) ValidateAccessToken(ctx context.Context, accessToken string) (string, *errx.Error) {
	t, err := s.verifyToken(accessToken)
	if err != nil {
		return "", err
	}

	if t.ExpiresAt.Before(time.Now()) {
		return "", errx.ErrToken
	}

	session, err := s.GetSession(ctx, t.SessionID)
	if err != nil {
		return "", err
	}

	if session.LastRefreshedAt.Equal(t.IssuedAt.Time) {
		return "", errx.ErrToken
	}

	return session.UserID, nil
}
