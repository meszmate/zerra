package token

import (
	"context"

	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/geo"
	"github.com/meszmate/zerra/internal/repostory"
)

type TokenService interface {
	GenerateSession(ctx context.Context, userID, ipaddr, userAgent string) (*models.Token, *errx.Error)
	GetSession(ctx context.Context, sessionID string) (*models.Session, *errx.Error)
	ValidateAccessToken(ctx context.Context, accessToken string) (string, *errx.Error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.Token, *errx.Error)
}

type tokenService struct {
	db             *db.DB
	tokenRepostory repostory.TokenRepostory
	geo            *geo.Client
	Cache          *cache.Cache

	AuthSecret string
}

func NewService(db *db.DB, cache *cache.Cache, geo *geo.Client, authSecret string) TokenService {
	return &tokenService{
		db:             db,
		tokenRepostory: repostory.NewTokenRepostory(db),
		geo:            geo,
		Cache:          cache,
		AuthSecret:     authSecret,
	}
}
