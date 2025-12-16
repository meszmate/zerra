package repostory

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
)

type UserRepostory interface {
	CreateUser(ctx context.Context, email, password string) (*models.User, *errx.Error)
	GetUser(ctx context.Context, id string) (*models.User, *errx.Error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, *errx.Error)
}

type userRepostory struct {
	DB *db.DB
}

func NewUserRepostory(db *db.DB) UserRepostory {
	return &userRepostory{
		DB: db,
	}
}

func (r *userRepostory) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, *errx.Error) {
	id := uuid.NewString()
	vMail, err := mail.ParseAddress(email)
	if err != nil {
		return nil, errx.ErrEmail
	}
	firstName := vMail.Name
	lastName := ""
	createdAt := time.Now()

	const q = `
		INSERT INTO users (id, email, password_hash, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
	`

	var params = []any{
		id,
		email,
		passwordHash,
		firstName,
		lastName,
		createdAt,
	}

	_, err = r.DB.Exec(
		ctx,
		q,
		params...)
	if err != nil {
		db.CaptureError(err, q, nil, "exec")
		return nil, errx.InternalError()
	}
	return &models.User{
		ID: id,

		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Roles:     make([]string, 0),

		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}, nil
}

func (r *userRepostory) getUser(ctx context.Context, key, value string) (*models.User, *errx.Error) {
	var u models.User

	q := fmt.Sprintf(
		`SELECT u.email, u.updated_at, u.created_at,
		   COALESCE(array_agg(ur.role_id) FILTER (WHERE ur.role_id IS NOT NULL), '{}') AS role_ids
		  FROM users u
		  LEFT JOIN user_roles ur ON ur.user_id = u.id
		  WHERE u.%s = $1
		  GROUP BY u.id`,
		key,
	)

	var params = []any{
		value,
	}

	err := r.DB.QueryRow(
		ctx,
		q,
		params...,
	).Scan(&u.Email, &u.UpdatedAt, &u.CreatedAt, &u.Roles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errx.ErrUser
		}
		db.CaptureError(err, q, params, "queryrow")
		return nil, errx.InternalError()
	}
	return &u, nil
}

func (r *userRepostory) GetUser(ctx context.Context, userID string) (*models.User, *errx.Error) {
	return r.getUser(ctx, "id", userID)
}

func (r *userRepostory) GetUserByEmail(ctx context.Context, email string) (*models.User, *errx.Error) {
	return r.getUser(ctx, "email", email)
}

func (r *userRepostory) UpdateUser(ctx context.Context, userID string, data *models.UpdateUser) (*models.User, *errx.Error) {
	// setClauses := []string{}
	// args := []any{userID}
	// argPos := 2
	return nil, nil
}
