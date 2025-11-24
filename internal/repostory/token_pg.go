package repostory

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
)

type TokenRepostory interface {
	GenerateSession(ctx context.Context, tx pgx.Tx, session *models.Session) *errx.Error
	GetSession(ctx context.Context, sessionID string) (*models.Session, *errx.Error)
	RefreshToken(ctx context.Context, sessionID, oldRefreshNonce, refreshNonce, accessNonce string, issuedAt time.Time) *errx.Error
	RevokeSession(ctx context.Context, tx pgx.Tx, sessionID string, revokedAt time.Time) *errx.Error

	FindExpiredSessions(ctx context.Context, userID string, cutoff time.Time) ([]string, *errx.Error)
	RevokeSessions(ctx context.Context, userID string) *errx.Error
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
		db.CaptureError(err, query, params, "exec")
		return errx.InternalError()
	}

	if cmd.RowsAffected() == 0 {
		return errx.ErrToken
	}

	return nil
}

func (r *tokenRepostory) RevokeSession(ctx context.Context, tx pgx.Tx, sessionID string, revokedAt time.Time) *errx.Error {
	query := `
		UPDATE sessions
		SET revoked_at = $1
		WHERE id = $2
	`

	params := []any{
		revokedAt, sessionID,
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

func (r *tokenRepostory) FindExpiredSessions(ctx context.Context, userID string, cutoff time.Time) ([]string, *errx.Error) {
	query := `
        SELECT id
        FROM sessions
        WHERE revoked_at IS NULL
          AND last_refreshed_at < $1
    `

	params := []any{
		cutoff,
	}

	rows, err := r.DB.Query(
		ctx,
		query,
		cutoff,
	)
	if err != nil {
		db.CaptureError(err, query, params, "query")
		return nil, errx.InternalError()
	}
	defer rows.Close()

	var sessions []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			db.CaptureError(err, "", nil, "scan")
			return nil, errx.InternalError()
		}
		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		db.CaptureError(err, "", nil, "rows_err")
		return nil, errx.InternalError()
	}

	return sessions, nil
}

func (r *tokenRepostory) RevokeSessions(ctx context.Context, userID string) *errx.Error {
	query := `
        UPDATE sessions
		SET revoked_at = now()
        WHERE revoked_at IS NULL
		  AND user_id = $1
    `

	params := []any{
		userID,
	}

	_, err := r.DB.Exec(
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
