package file

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/redis/go-redis/v9"
)

func defaultAvatarsKey() string {
	return "default_avatars"
}

func (s *fileService) saveDefaultAvatars(ctx context.Context, avatars []string) *errx.Error {
	key := defaultAvatarsKey()

	data, err := json.Marshal(avatars)
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if err := s.cache.SetEx(ctx, key, data, DefaultAvatarsTTL).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *fileService) getDefaultAvatars(ctx context.Context) ([]string, *errx.Error) {
	key := defaultAvatarsKey()

	data, err := s.cache.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	var avatars []string

	if err := json.Unmarshal(data, &avatars); err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	return avatars, nil
}
