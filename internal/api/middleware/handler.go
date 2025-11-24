package middleware

import "github.com/meszmate/zerra/internal/app/token"

type Handler struct {
	tokenService token.TokenService
}
