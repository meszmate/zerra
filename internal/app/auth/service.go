package auth

import (
	"context"

	"github.com/meszmate/zerra/internal/app/token"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/captcha"
	"github.com/meszmate/zerra/internal/repostory"
)

type AuthService interface {
	LoginStart(ctx context.Context, data *models.LoginStart) (*models.AuthSession, *errx.Error)
}

type authService struct {
	db            *db.DB
	authRepostory repostory.AuthRepostory
	tokenService  token.TokenService
	Cache         *cache.Cache
	captcha       *captcha.Turnstile
}

func NewService(db *db.DB, cache *cache.Cache, tokenService token.TokenService, captcha *captcha.Turnstile) AuthService {
	return &authService{
		db:            db,
		authRepostory: repostory.NewAuthRepostory(db),
		tokenService:  tokenService,
		Cache:         cache,
		captcha:       captcha,
	}
}
