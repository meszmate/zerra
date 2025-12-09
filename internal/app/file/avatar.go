package file

import (
	"context"

	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/pkg/crypt"
)

func (s *fileService) GetDefaultAvatars(ctx context.Context, userID string) (string, *errx.Error) {
	avatars, err := s.getDefaultAvatars(ctx)
	if err != nil {
		return "", err
	}

	avatars, err = s.fileRepostory.GetDefaultAvatars(ctx)
	if err != nil {
		return "", err
	}

	if err := s.saveDefaultAvatars(ctx, avatars); err != nil {
		return "", err
	}

	return avatars[crypt.HashToRange(userID, len(avatars)-1)], nil
}

func (s *fileService) UploadAvatar(ctx context.Context, userID string) (string, *errx.Error) {
	s.fileService

}
