package repostory

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/gen"
)

type TokenRepostory interface {
	GenerateSession(ctx context.Context, tx pgx.Tx, session *models.Session) *errx.Error
	GetSession(ctx context.Context, sessionID string) (*models.Session, *errx.Error)
	RefreshToken(ctx context.Context, sessionID, oldRefreshNonce, refreshNonce, accessNonce string) *errx.Error
}

type tokenRepostory struct {
	DB *db.DB
}

func NewTokenRepostory(db *db.DB) TokenRepostory {
	return &tokenRepostory{
		DB: db,
	}
}

func (r *tokenRepostory) GenerateSession(ctx context.Context, tx pgx.Tx, session *models.Session) *errx.Error {
	query := `
		INSERT INTO sessions (
		 id, user_id,
		 created_at, expires_at, last_refreshed_at, revoked_at,
		 access_nonce, refresh_nonce,
		 location_city, location_region, location_country, location_country_code, location_postal_code,
		 os_name, browser_name
		) VALUES (
		 $1, $2,
		 $3, $4, $5, $6,
		 $7, $8,
		 $9, $10, $11, $12, $13,
		 $14, $15
		)
	`

	params := []any{
		session.ID, session.UserID,
		session.CreatedAt, session.ExpiresAt, session.LastRefreshedAt, session.RevokedAt,
		session.AccessNonce, session.RefreshNonce,
		session.LocationCity, session.LocationRegion, session.LocationCountry, session.LocationCountryCode, session.LocationPostalCode,
		session.OSName, session.BrowserName,
	}

	_, err := tx.Exec(
		ctx,
		query,
		params...,
	)
	if err != nil {
		db.CaptureError(err, query, params, "exec")
		return errx.InternalError()
	}
	return nil
}

func (r *tokenRepostory) GetSession(ctx context.Context, sessionID string) (*models.Session, *errx.Error) {
	query := `
		SELECT
		 id, user_id,
		 created_at, expires_at, last_refreshed_at, revoked_at,
		 access_nonce, refresh_nonce,
		 location_city, location_region, location_country, location_country_code, location_postal_code,
		 os_name, browser_name
		FROM sessions
		WHERE id = $1
	`

	params := []any{
		sessionID,
	}

	var sess models.Session

	err := r.DB.QueryRow(
		ctx,
		query,
		params...,
	).Scan(
		&sess.ID, &sess.UserID,
		&sess.CreatedAt, &sess.ExpiresAt, &sess.LastRefreshedAt, &sess.RevokedAt,
		&sess.AccessNonce, &sess.RefreshNonce,
		&sess.LocationCity, &sess.LocationRegion, &sess.LocationCountry, &sess.LocationCountryCode, &sess.LocationPostalCode,
		&sess.OSName, &sess.BrowserName,
	)
	if err != nil {
		db.CaptureError(err, query, params, "queryrow")
		return nil, errx.InternalError()
	}

	return &sess, nil
}

func (r *tokenRepostory) RefreshToken(ctx context.Context, sessionID, oldRefreshNonce, accessNonce, refreshNonce string, issuedAt time.Time) *errx.Error {
	query := `
		UPDATE sessions
		SET last_refreshed_at = $5,
		 access_nonce = $1, refresh_nonce = $2
		WHERE refresh_nonce = $3 AND id = $4
	`

	params := []any{
		accessNonce, refreshNonce,
		oldRefreshNonce, sessionID,
		issuedAt,
	}

	cmd, err := r.DB.Exec(
		ctx,
		query,
		params...,
	)
	if err != nil {
		db.CaptureError(err, query, params, "queryrow")
		return errx.InternalError()
	}

	if cmd.RowsAffected() == 0 {
		return errx.ErrToken
	}

	return nil
}

func (r *tokenRepostory) GenerateEmailVerificationCode(ctx context.Context, email string) (string, *errx.Error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		db.CaptureError(err, "", nil, "begin")
		return "", errx.InternalError()
	}
	defer tx.Rollback(ctx)

	code, err := gen.VerificationCode()
	if err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}
	encryptedCode, err := r.Encrypt.Encrypt(code)
	if err != nil {
		sentry.CaptureException(err)
		return "", errx.InternalError()
	}
	expiresAt := time.Now().Add(time.Duration(config.EMAIL_VERIFICATION_EXPIRATION))

	query := `
		INSERT INTO email_verification (email, code, expires_at)
		VALUES ($1, $2, $3)
	`

	params := []any{
		email,
		encryptedCode,
		expiresAt,
	}

	_, err = tx.Exec(
		ctx,
		query,
		params...,
	)
	if err != nil {
		db.CaptureError(err, query, params, "exec")
		return "", errx.InternalError()
	}
	if err := tx.Commit(ctx); err != nil {
		db.CaptureError(err, "", nil, "commit")
		return "", errx.InternalError()
	}

	return code, nil
}

func (r *tokenRepostory) GeneratePasswordVerificationToken(ctx context.Context, email string) (string, *errx.Error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		db.CaptureError(err, "", nil, "begin")
		return "", errx.InternalError()
	}
	defer tx.Rollback(ctx)

	params := []any{
		email,
	}

	var userid string
	var exists bool
	err = tx.QueryRow(
		ctx,
		query,
		params...,
	).Scan(&userid, &exists)
	if err != nil {
		db.CaptureError(err, query, params, "queryrow")
		return "", errx.InternalError()
	}

	if exists {
		return "", errx.Userf("There is already an active password reset token.")
	}

	token, err := gen.Token(gen.AUTH_SESSION_TOKEN_LEN)
	if err != nil {
		return "", errx.Internalf(db.Cfg.Log, "DB: GeneratePasswordVerificationToken token generation", err)
	}

	expiresAt := time.Now().Add(time.Duration(config.RESET_PASSWORD_EXPIRATION))

	_, err = tx.Exec(ctx,
		"INSERT INTO password_reset (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userid, crypto.HashTokenSHA256(token), expiresAt)
	if err != nil {
		return "", errx.Internalf(db.Cfg.Log, "DB: GeneratePasswordVerificationToken exec", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return "", errx.Internalf(db.Cfg.Log, "DB: GeneratePasswordVerificationToken commit", err)
	}

	return token, nil
}
