package handler

import (
	"github.com/meszmate/zerra/internal/app/auth"
	"github.com/meszmate/zerra/internal/app/token"
)

type Handler struct {
	tokenService token.TokenService
	authService  auth.AuthService
}
