package file

import (
	"context"
	"io"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
)

func (s *fileService) PutObject(ctx context.Context, fileOwnerType FileOwnerType, contentType, key string, body io.Reader) *errx.Error {
	s.fileRepostory.

	if err := s.client.PutObject(ctx, key, contentType, body); err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}
