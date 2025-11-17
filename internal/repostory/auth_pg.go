package repostory

import (
	"context"
	"net/mail"
	"os/user"
	"time"

	"github.com/google/uuid"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/pkg/argon2"
)

type AuthRepostory interface {
}

type authRepostory struct {
	DB *db.DB
}

func NewAuthRepostory(db *db.DB) AuthRepostory {
	return &authRepostory{
		DB: db,
	}
}

func (r *authRepostory) IsValidCredentials(ctx context.Context, email, password string) bool {
	var pw string
	err := r.DB.QueryRow(ctx, "SELECT password_hash FROM users WHERE email = $1", email).Scan(&pw)
	if err != nil {
		return false
	}
	return argon2.VerifyPassword(password, pw)
}

func (r *authRepostory) CleanCode(ctx context.Context, email string) *errx.Error {
	q := `DELETE FROM email_verification
		WHERE email = $1`

	var params = []any{
		email,
	}

	_, err := r.DB.Exec(
		ctx,
		q,
		params...,
	)

	if err != nil {
		db.CaptureError(err, q, nil, "exec")
		return errx.InternalError()
	}

	return nil
}

func (r *userRepostory) ExternalLogin(ctx context.Context, email string) (*user.User, *errx.Error) {
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

	var u user.User
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
