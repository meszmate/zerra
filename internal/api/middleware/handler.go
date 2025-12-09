package middleware

import "github.com/meszmate/zerra/internal/app/token"

type Handler struct {
	TokenService token.TokenService
}
