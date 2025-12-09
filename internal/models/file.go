package models

import "time"

type FileParentType string

const (
	FileParentTypeAvatar FileParentType = "avatar"
	FileParentTypeIcon   FileParentType = "icon"
)

type File struct {
	ID         string         `json:"id"`
	ParentType FileParentType `json:"parent_type"`
	ParentID   string         `json:"parent_id"`
	Name       string         `json:"name"`
	FileType   string         `json:"file_type"`
	CreatedAt  time.Time      `json:"created_at"`
}
