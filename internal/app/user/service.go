package user

import (
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/repostory"
)

type UserService interface {
}

type userService struct {
	db            *db.DB
	Cache         *cache.Cache
	userRepostory repostory.UserRepostory
}

func NewService(db *db.DB, cache *cache.Cache) UserService {
	return &userService{
		db:    db,
		Cache: cache,
	}
}
