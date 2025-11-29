package models

import "time"

type Organization struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`

	// Optional
	StudentsCount int `json:"students_count"`
	TeachersCount int `json:"teachers_count"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
