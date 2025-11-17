package token

import (
	"context"
	"net/netip"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meszmate/zerra/internal/errx"
	"github.com/meszmate/zerra/internal/infrastructure/db"
	"github.com/meszmate/zerra/internal/models"
	"github.com/meszmate/zerra/internal/pkg/crypt"
	"github.com/mileusna/useragent"
)

type TokenClaims struct {
	UserID    string `json:"sub"`
	SessionID string `json:"sid"`
	Nonce     string `json:"nonce"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a short-lived JWT.
func (s *tokenService) generateToken(userID, sessionID, nonce string, issuedAt, expiresAt time.Time) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		Nonce:     nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.AuthSecret)
}

func (s *tokenService) verifyToken(tokenStr string) (*TokenClaims, *errx.Error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errx.ErrToken
		}
		return s.AuthSecret, nil
	})

	if err != nil {
		return nil, errx.ErrToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errx.ErrToken
	}

	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return nil, errx.ErrToken
	}

	return claims, nil
}

func (s *tokenService) GenerateSession(ctx context.Context, userID string, ipaddr string, userAgent string) (*models.Token, *errx.Error) {
	ip, err := netip.ParseAddr(ipaddr)
	if err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	ipinfo, err := s.geo.Lookup(ip)
	if err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	userAgentInfo := useragent.Parse(userAgent)

	session := &models.Session{
		ID:     uuid.NewString(),
		UserID: userID,

		LocationCity:       ipinfo.City,
		LocationRegion:     ipinfo.Region,
		LocationCountry:    ipinfo.Country,
		LocationPostalCode: ipinfo.PostalCode,

		BrowserName: userAgentInfo.Name,
		OSName:      userAgentInfo.Name,
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		db.CaptureError(err, "", nil, "begin")
		return nil, errx.InternalError()
	}
	defer tx.Rollback(ctx)

	issuedAt := time.Now()
	session.LastRefreshedAt = issuedAt
	session.CreatedAt = issuedAt

	accessTokenExpiresAt := issuedAt.Add(10 * time.Minute)
	accessNonce, err := crypt.Nonce()
	if err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}
	session.AccessNonce = accessNonce

	accessToken, err := s.generateToken(userID, session.ID, accessNonce, issuedAt, accessTokenExpiresAt)
	if err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	refreshTokenExpiresAt := issuedAt.Add(2 * 30 * 24 * time.Hour)
	refreshNonce, err := crypt.Nonce()
	if err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}
	session.RefreshNonce = refreshNonce

	refreshToken, err := s.generateToken(userID, session.ID, refreshNonce, issuedAt, refreshTokenExpiresAt)
	if err != nil {
		sentry.CaptureException(err)
		return nil, errx.InternalError()
	}

	if err := s.tokenRepostory.GenerateSession(ctx, tx, session); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		db.CaptureError(err, "", nil, "commit")
		return nil, errx.InternalError()
	}

	return &models.Token{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}
