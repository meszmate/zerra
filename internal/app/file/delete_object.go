package file

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
)

func (s *fileService) DeleteObject(ctx context.Context, id string) *errx.Error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	f, xerr := s.fileRepostory.DeleteFile(ctx, tx, id)
	if xerr != nil {
		return xerr
	}

	if xerr := s.client.DeleteObject(ctx, GetKey(f.ParentType, f.ParentID, f.Name)); xerr != nil {
		return xerr
	}

	if err := tx.Commit(ctx); err != nil {
		return errx.InternalError()
	}

	return nil
}
