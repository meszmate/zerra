package config

import "os"

type Config struct {
	PrimaryDB string
	Redis     string
	GeoDBPath string
	S3Bucket  string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectUri  string

	NotifyName    string
	NotifyAddress string

	Hostname string
	Port     string

	SentryDSN string

	AuthSecret      string
	TurnstileSecret string
}

func Load() *Config {
	return &Config{
		PrimaryDB: os.Getenv("PRIMARY_DB"),
		Redis:     os.Getenv("REDIS"),
		GeoDBPath: os.Getenv("GEO_DB_PATH"),
		S3Bucket:  os.Getenv("S3_BUCKET"),

		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectUri:  os.Getenv("GOOGLE_REDIRECT_URI"),

		NotifyName:    os.Getenv("NOTIFY_NAME"),
		NotifyAddress: os.Getenv("NOTIFY_ADDRESS"),

		Hostname: os.Getenv("HOSTNAME"),
		Port:     os.Getenv("PORT"),

		SentryDSN: os.Getenv("SENTRY_DSN"),

		AuthSecret:      os.Getenv("AUTH_SECRET"),
		TurnstileSecret: os.Getenv("TURNSTILE_SECRET"),
	}
}
