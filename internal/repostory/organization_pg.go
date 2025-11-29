package repostory

import (
	"context"

	"github.com/meszmate/zerra/internal/infrastructure/db"
)

type OrganizationRepostory interface {
}

type organizationRepostory struct {
	DB *db.DB
}

func NewOrganizationRepostory(db *db.DB) AuthRepostory {
	return &authRepostory{
		DB: db,
	}
}

func (r *organizationRepostory) Get(ctx context.Context, userID string) 
