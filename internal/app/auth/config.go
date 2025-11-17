package auth

import "time"

const (
	AuthSessionTTL = 10 * time.Minute
	AuthEmailTTL   = 30 * time.Minute
	AuthEmailLimit = 5
)
