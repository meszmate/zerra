package organization

import (
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/infrastructure/db"
)

type OrganizationService interface {
}

type organizationService struct {
	db    *db.DB
	cache *cache.Cache
}

func NewService(db *db.DB, cache *cache.Cache) OrganizationService {
	return &organizationService{
		db:    db,
		cache: cache,
	}
}
