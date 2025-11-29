package config

import "os"

type Config struct {
	AppDomain string
}

func Load() *Config {
	return &Config{
		AppDomain: os.Getenv("APP_DOMAIN"),
	}
}
