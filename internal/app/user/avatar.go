package user

import (
	"context"
	"mime/multipart"

	"github.com/meszmate/zerra/internal/config"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/utils"
)

func (s *userService) ChangeAvatar(ctx context.Context, userID string, avatar *multipart.FileHeader) (string, *errx.Error) {
	imgreader, err := utils.GetJPG(avatar, config.MaxAvatarSizeBytes, config.AvatarSize, config.AvatarSize)
	if err != nil {
		return "", err
	}

	s.file
}
