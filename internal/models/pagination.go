package models

type Pagination struct {
	Total      *int64  `json:"total"`
	NextCursor *string `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
}
