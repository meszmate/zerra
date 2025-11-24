package user

import (
	"context"

	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/repostory"
)

type UserService interface {
	SaveUser(ctx context.Context, user *models.User) *errx.Error
	GetUser(ctx context.Context, userID string) (*models.User, *errx.Error)
}

type userService struct {
	db            *db.DB
	cache         *cache.Cache
	userRepostory repostory.UserRepostory
}

func NewService(db *db.DB, cache *cache.Cache) UserService {
	return &userService{
		db:    db,
		cache: cache,
	}
}
