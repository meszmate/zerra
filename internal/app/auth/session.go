package auth

import (
	"context"

	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
)

func (s *authService) LoginStart(ctx context.Context, data *models.LoginStart, ipAddr string) (*models.AuthSession, *errx.Error) {
	c, err := s.captcha.Verify(ctx, data.Turnstile, ipAddr)
	if err != nil {
		return nil, err
	}

	s.authRepostory.
}
