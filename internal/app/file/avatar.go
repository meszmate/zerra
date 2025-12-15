package file

import (
	"context"

	"github.com/google/uuid"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/crypt"
)

func (s *fileService) GetDefaultAvatar(ctx context.Context, userID string) (*models.File, *errx.Error) {
	avatars, err := s.getDefaultAvatars(ctx)
	if err != nil {
		return nil, err
	}

	avatars, err = s.fileRepostory.GetFilesByParent(ctx, models.FileParentTypeAvatar, uuid.Nil.String())
	if err != nil {
		return nil, err
	}

	if err := s.saveDefaultAvatars(ctx, avatars); err != nil {
		return nil, err
	}

	return &avatars[crypt.HashToRange(userID, len(avatars)-1)], nil
}

func (s *fileService) UploadAvatar(ctx context.Context, userID string) (string, *errx.Error) {
	return "", nil
}
