package auth

import (
	"context"

	"github.com/meszmate/zerra/internal/app/token"
	"github.com/meszmate/zerra/internal/app/user"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/notify"
	"github.com/meszmate/zerra/internal/pkg/captcha"
	"github.com/meszmate/zerra/internal/repostory"
)

type AuthService interface {
	LoginStart(ctx context.Context, data *AuthData, ipaddr string) (*models.AuthSession, *errx.Error)
	LoginConfirm(ctx context.Context, data *ConfirmData, session, ipaddr, userAgent string) (*models.Token, *errx.Error)

	RegistrationStart(ctx context.Context, data *AuthData, ipaddr string) (*models.AuthSession, *errx.Error)
	RegistrationConfirm(ctx context.Context, data *ConfirmData, session, ipaddr string) *errx.Error

	ResetPasswordStart(ctx context.Context, data *ResetPasswordStart, ipaddr string) *errx.Error
	ResetPasswordConfirm(ctx context.Context, data *ResetPasswordConfirm, session, ipaddr string) *errx.Error
}

type authService struct {
	authRepostory            repostory.AuthRepostory
	userRepostory            repostory.UserRepostory
	tokenService             token.TokenService
	userService              user.UserService
	emailNotificationService notify.EmailNotificationService

	cache   *cache.Cache
	captcha *captcha.Turnstile
}

func NewService(
	authRepostory repostory.AuthRepostory,
	userRepostory repostory.UserRepostory,
	tokenService token.TokenService,
	userService user.UserService,
	emailNotificationService notify.EmailNotificationService,
	cache *cache.Cache,
	captcha *captcha.Turnstile,
) AuthService {
	return &authService{
		authRepostory:            authRepostory,
		userRepostory:            userRepostory,
		tokenService:             tokenService,
		userService:              userService,
		emailNotificationService: emailNotificationService,

		cache:   cache,
		captcha: captcha,
	}
}
