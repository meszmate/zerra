package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
	"github.com/redis/go-redis/v9"
)

func getUserKey(id string) string {
	return "user:" + id
}

func (s *userService) SaveUser(ctx context.Context, user *models.User) *errx.Error {
	key := getUserKey(user.ID)

	data, err := json.Marshal(user)
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if err := s.cache.SetEx(ctx, key, data, UserTTL).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *userService) getUser(ctx context.Context, userID string) (*models.User, *errx.Error) {
	data, err := s.cache.Get(ctx, getUserKey(userID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	return &user, nil
}
