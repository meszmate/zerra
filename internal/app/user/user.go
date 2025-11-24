package user

import (
	"context"

	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
)

func (s *userService) GetUser(ctx context.Context, userID string) (*models.User, *errx.Error) {
	u, err := s.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return u, nil
	}

	u, err = s.userRepostory.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.SaveUser(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
