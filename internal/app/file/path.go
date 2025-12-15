package file

import (
	"github.com/meszmate/zerra/internal/models"
)

var (
	AvatarPath = func(userID string, name string) string {
		return "avatar/" + userID + "/" + name
	}
	IconPath = func(organizationID string, name string) string {
		return "icons/" + organizationID + "/" + name
	}
)

func GetKey(fileParentType models.FileParentType, fileParentID string, name string) string {
	var key string
	switch fileParentType {
	case models.FileParentTypeAvatar:
		key = AvatarPath(fileParentID, name)
	case models.FileParentTypeIcon:
		key = IconPath(fileParentID, name)
	}

	return key
}
