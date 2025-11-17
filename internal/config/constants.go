package config

import "time"

const (
	DefaultColor = "#c4c8cf"
	Domain       = "warmbly.com"

	EmailVerificationTTL = 4 * time.Hour
	ResetPasswordTTL     = 4 * time.Hour
	AccessTokenTTL       = 15 * time.Minute
	RefreshTokenTTL      = 60 * 24 * time.Hour
)
