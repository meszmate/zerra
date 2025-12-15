package file

import (
	"context"
	"io"

	"github.com/jackc/pgx/v5"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/models"
)

func (s *fileService) PutObject(ctx context.Context, tx pgx.Tx, fileParentType models.FileParentType, fileParentID, contentType string, name string, body io.Reader) (*models.File, *errx.Error) {
	var key string
	switch fileParentType {
	case models.FileParentTypeAvatar:
		key = AvatarPath(fileParentID, name)
	case models.FileParentTypeIcon:
		key = IconPath(fileParentID, name)
	}

	file, xerr := s.fileRepostory.PutFile(ctx, tx, key, contentType, models.FileParentType(fileParentType), fileParentID, name)
	if xerr != nil {
		return nil, xerr
	}

	if err := s.client.PutObject(ctx, key, contentType, body); err != nil {
		return nil, err
	}

	return file, nil
}
