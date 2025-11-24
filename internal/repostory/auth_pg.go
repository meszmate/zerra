package repostory

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/argon2"
)

type AuthRepostory interface {
	IsValidCredentials(ctx context.Context, email, password string) (string, *errx.Error)
	ExternalLogin(ctx context.Context, email string) (*models.User, *errx.Error)
	ResetPassword(ctx context.Context, userID string, password string) *errx.Error
}

type authRepostory struct {
	DB *db.DB
}

func NewAuthRepostory(db *db.DB) AuthRepostory {
	return &authRepostory{
		DB: db,
	}
}

func (r *authRepostory) IsValidCredentials(ctx context.Context, email, password string) (string, *errx.Error) {
	var id string
	var pw string

	query := `
		SELECT id, password_hash
		FROM users
		WHERE email = $1
	`

	params := []any{
		email,
	}

	err := r.DB.QueryRow(
		ctx,
		query,
		params...,
	).Scan(&id, &pw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errx.ErrCredentials
		}
		db.CaptureError(err, query, params, "queryrow")
		return "", errx.InternalError()
	}

	val, err := argon2.Verify(password, pw)
	if err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}

	if !val {
		return "", errx.ErrCredentials
	}

	return id, nil
}

func (r *authRepostory) ExternalLogin(ctx context.Context, email string) (*models.User, *errx.Error) {
	id := uuid.NewString()
	vMail, err := mail.ParseAddress(email)
	if err != nil {
		return nil, errx.ErrEmail
	}

	firstName := vMail.Name
	lastName := ""
	now := time.Now()

	query := `
        INSERT INTO users (id, email, password_hash, first_name, last_name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $6)
        ON CONFLICT (email) DO UPDATE
        SET updated_at = EXCLUDED.updated_at
        RETURNING id, email, first_name, last_name, created_at, updated_at;
    `

	var params = []any{
		id,
		email,
		"",
		firstName,
		lastName,
		now,
	}

	var u models.User
	err = r.DB.QueryRow(
		ctx,
		query,
		params...,
	).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		db.CaptureError(err, query, nil, "queryrow")
		return nil, errx.InternalError()
	}

	return &u, nil
}

func (r *authRepostory) ResetPassword(ctx context.Context, userID string, passwordHash string) *errx.Error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = now()
		WHERE id = $2
	`
	params := []any{
		passwordHash,
		userID,
	}

	if _, err := r.DB.Exec(ctx, query, params...); err != nil {
		db.CaptureError(err, query, params, "exec")
		return errx.InternalError()
	}

	return nil
}
