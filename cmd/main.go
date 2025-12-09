package main

import (
	"context"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/api"
	"github.com/meszmate/zerra/internal/api/handler"
	"github.com/meszmate/zerra/internal/api/middleware"
	"github.com/meszmate/zerra/internal/app/auth"
	"github.com/meszmate/zerra/internal/app/file"
	"github.com/meszmate/zerra/internal/app/organization"
	"github.com/meszmate/zerra/internal/app/token"
	"github.com/meszmate/zerra/internal/app/user"
	"github.com/meszmate/zerra/internal/config"
	"github.com/meszmate/zerra/internal/infrastructure/cache"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/notify"
	"github.com/meszmate/zerra/internal/pkg/captcha"
	"github.com/meszmate/zerra/internal/pkg/geo"
	"github.com/meszmate/zerra/internal/repostory"
)

func main() {
	var addr string

	var authService auth.AuthService
	var tokenService token.TokenService
	var userService user.UserService
	var organizationService organization.OrganizationService
	var fileService file.FileService
	{
		cfg := config.Load()

		if err := sentry.Init(sentry.ClientOptions{
			Dsn:            cfg.SentryDSN,
			SendDefaultPII: true,
		}); err != nil {
			log.Fatal(err)
		}

		db, err := db.New(context.Background(), cfg.PrimaryDB)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}

		cache, err := cache.New(cfg.Redis)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}

		geoClient, err := geo.New(cfg.GeoDBPath)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}

		emailNotificationService, err := notify.NewEmailNotficiationService(context.Background(), cfg.NotifyName, cfg.NotifyAddress)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}

		captcha := captcha.NewTurnstile(cfg.TurnstileSecret)

		fileClient, err := file.NewClient(context.Background(), cfg.S3Bucket)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}

		authRepostory := repostory.NewAuthRepostory(db)
		tokenRepostory := repostory.NewTokenRepostory(db)
		userRepostory := repostory.NewUserRepostory(db)
		organizationRepostory := repostory.NewOrganizationRepostory(db)
		fileRepostory := repostory.NewFileRepostory(db)

		tokenService = token.NewService(db, tokenRepostory, cache, geoClient, cfg.AuthSecret)
		authService = auth.NewService(
			authRepostory,
			userRepostory,
			tokenService,
			userService,
			emailNotificationService,
			cache,
			captcha,
		)
		fileService = file.NewService(fileClient, fileRepostory, userRepostory, cache)
		userService = user.NewService(userRepostory, fileService, cache)
		organizationService = organization.NewService(organizationRepostory, cache)

		addr = cfg.Hostname + ":" + cfg.Port
	}

	h := &handler.Handler{
		AuthService:         authService,
		TokenService:        tokenService,
		UserService:         userService,
		OrganizationService: organizationService,
		FileService:         fileService,
	}

	m := &middleware.Handler{
		TokenService: tokenService,
	}

	api.Run(h, m, addr)
}
