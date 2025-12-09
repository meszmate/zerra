package handler

import (
	"github.com/meszmate/zerra/internal/app/auth"
	"github.com/meszmate/zerra/internal/app/file"
	"github.com/meszmate/zerra/internal/app/organization"
	"github.com/meszmate/zerra/internal/app/token"
	"github.com/meszmate/zerra/internal/app/user"
)

type Handler struct {
	AuthService         auth.AuthService
	TokenService        token.TokenService
	UserService         user.UserService
	OrganizationService organization.OrganizationService
	FileService         file.FileService
}
