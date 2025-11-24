package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/meszmate/zerra/internal/api/handler"
	"github.com/meszmate/zerra/internal/api/middleware"
)

func New(h *handler.Handler, m *middleware.Handler, open bool, port string) {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PATCH", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/auth")
	{
		r.POST("/login/start", h.LoginStart)
		r.POST("/login/confirm", h.LoginConfirm)
		r.POST("/register/start", h.RegistrationStart)
		r.POST("/register/confirm", h.RegistrationConfirm)
		r.POST("/refresh", h.RefreshToken)
		r.POST("/reset-password/start", h.ResetPasswordStart)
		r.POST("/reset-password/confirm", h.ResetPasswordStart)
	}

	protectedAuth := auth.Group("")
	protectedAuth.Use(m.AuthMiddleware())
	{
		r.POST("/logout", h.Logout)
		r.POST("/logout-all", h.LogoutAll)
		r.GET("/me", h.GetUser)
	}
}
