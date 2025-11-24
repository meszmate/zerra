package handler

import (
	"github.com/meszmate/zerra/internal/app/auth"
	"github.com/meszmate/zerra/internal/app/token"
	"github.com/meszmate/zerra/internal/app/user"
)

type Handler struct {
	authService  auth.AuthService
	tokenService token.TokenService
	userService  user.UserService
}
