package file

import (
	"context"
	"io"

	"github.com/jackc/pgx/v5"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/repostory"
)

type FileService interface {
	GetDefaultAvatar(ctx context.Context, userID string) (*models.File, *errx.Error)
	PutObject(ctx context.Context, pgx pgx.Tx, fileParentType models.FileParentType, fileParentID, contentType, name string, body io.Reader) (*models.File, *errx.Error)
	DeleteObject(ctx context.Context, id string) *errx.Error
}

type fileService struct {
	db            *db.DB
	client        *Client
	cache         *cache.Cache
	userRepostory repostory.UserRepostory
	fileRepostory repostory.FileRepostory
}

func NewService(db *db.DB, client *Client, fileRepostory repostory.FileRepostory, userRepostory repostory.UserRepostory, cache *cache.Cache) FileService {
	return &fileService{
		db:            db,
		client:        client,
		userRepostory: userRepostory,
		fileRepostory: fileRepostory,
		cache:         cache,
	}
}
