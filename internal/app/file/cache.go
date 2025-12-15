package file

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
	"github.com/redis/go-redis/v9"
)

var defaultAvatarsKey = "default_avatars"

func (s *fileService) saveDefaultAvatars(ctx context.Context, avatars []models.File) *errx.Error {
	return s.saveFiles(ctx, defaultAvatarsKey, avatars, DefaultAvatarsTTL)
}

func (s *fileService) getDefaultAvatars(ctx context.Context) ([]models.File, *errx.Error) {
	return s.getFiles(ctx, defaultAvatarsKey)
}

func (s *fileService) saveFiles(ctx context.Context, key string, files []models.File, ttl time.Duration) *errx.Error {
	data, err := json.Marshal(files)
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	if err := s.cache.SetEx(ctx, key, data, ttl).Err(); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (s *fileService) getFiles(ctx context.Context, key string) ([]models.File, *errx.Error) {
	data, err := s.cache.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	var files []models.File

	if err := json.Unmarshal(data, &files); err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	return files, nil
}
