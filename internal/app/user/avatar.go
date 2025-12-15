package user

import (
	"context"
	"mime/multipart"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/config"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/utils"
)

func (s *userService) ChangeAvatar(ctx context.Context, userID string, avatar *multipart.FileHeader) (*models.User, *errx.Error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	var previousAvatar string

	// Lock the row
	query := `
		SELECT avatar_id
		FROM users
		WHERE id = $1
		FOR UPDATE
	`

	params := []any{
		userID,
	}

	if err := tx.QueryRow(ctx, query, params...).Scan(&previousAvatar); err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	if previousAvatar != "" {
		if err := s.fileService.DeleteObject(ctx, previousAvatar); err != nil {
			return nil, err
		}
	}

	imgreader, xerr := utils.GetJPG(avatar, config.MaxAvatarSizeBytes, config.AvatarSize, config.AvatarSize)
	if xerr != nil {
		return nil, xerr
	}

	name, xerr := utils.GetFileHash(imgreader)
	if err != nil {
		return nil, xerr
	}

	fullName := name + ".jpg"

	u, xerr := s.GetUser(ctx, userID)
	if xerr != nil {
		return nil, xerr
	}

	f, xerr := s.fileService.PutObject(ctx, tx, models.FileParentTypeAvatar, userID, "image/jpeg", fullName, imgreader)
	if xerr != nil {
		return nil, xerr
	}

	u.Avatar = f.Name

	s.SaveUser(ctx, u)

	if err := tx.Commit(ctx); err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	return u, nil
}
