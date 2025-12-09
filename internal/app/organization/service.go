package organization

import (
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/repostory"
)

type OrganizationService interface {
}

type organizationService struct {
	organizationRepostory repostory.OrganizationRepostory
	cache                 *cache.Cache
}

func NewService(organizationRepostory repostory.OrganizationRepostory, cache *cache.Cache) OrganizationService {
	return &organizationService{
		organizationRepostory: organizationRepostory,
		cache:                 cache,
	}
}
