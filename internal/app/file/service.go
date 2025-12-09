package file

import (
	"context"

	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/repostory"
)

type FileService interface {
	GetDefaultAvatars(ctx context.Context, userID string) (string, *errx.Error)
}

type fileService struct {
	client        *Client
	cache         *cache.Cache
	userRepostory repostory.UserRepostory
	fileRepostory repostory.FileRepostory
}

func NewService(client *Client, fileRepostory repostory.FileRepostory, userRepostory repostory.UserRepostory, cache *cache.Cache) FileService {
	return &fileService{
		client:        client,
		userRepostory: userRepostory,
		fileRepostory: fileRepostory,
		cache:         cache,
	}
}
