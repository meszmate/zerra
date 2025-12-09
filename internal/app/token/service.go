package token

import (
	"context"
	"time"

	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/geo"
	"github.com/meszmate/zerra/internal/repostory"
)

type TokenService interface {
	GenerateToken(userID, sessionID, nonce string, issuedAt, expiresAt time.Time) (string, error)
	VerifyToken(tokenStr string) (*TokenClaims, *errx.Error)
	GenerateSession(ctx context.Context, userID, ipaddr, userAgent string) (*models.Token, *errx.Error)
	GetSession(ctx context.Context, sessionID string) (*models.Session, *errx.Error)
	ValidateAccessToken(ctx context.Context, accessToken string) (*models.Session, *errx.Error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.Token, *errx.Error)

	RevokeSession(ctx context.Context, accessToken string) *errx.Error
	RevokeAllSession(ctx context.Context, accessToken string) *errx.Error
}

type tokenService struct {
	db             *db.DB
	tokenRepostory repostory.TokenRepostory
	geo            *geo.Client
	cache          *cache.Cache

	AuthSecret string
}

func NewService(db *db.DB, tokenRepostory repostory.TokenRepostory, cache *cache.Cache, geo *geo.Client, authSecret string) TokenService {
	return &tokenService{
		db:             db,
		tokenRepostory: tokenRepostory,
		geo:            geo,
		cache:          cache,

		AuthSecret: authSecret,
	}
}
