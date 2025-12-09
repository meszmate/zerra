package file

import "time"

const (
	DefaultAvatarsTTL = 2 * time.Hour
)

type FileOwnerType string

const (
	FileOwnerTypeUser         FileOwnerType = "user"
	FileOwnerTypeOrganization FileOwnerType = "organization"
)
