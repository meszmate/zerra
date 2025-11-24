package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meszmate/zerra/internal/api/middleware"
	"github.com/meszmate/zerra/internal/app/auth"
	"github.com/meszmate/zerra/internal/errx"
)

func (h *Handler) LoginStart(c *gin.Context) {
	var data auth.AuthData

	if err := c.ShouldBindJSON(&data); err != nil {
		errx.Handle(c, errx.ErrInvalid)
		return
	}

	resp, err := h.authService.LoginStart(c.Request.Context(), &data, c.ClientIP())
	if err != nil {
		errx.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) LoginConfirm(c *gin.Context) {
	var data auth.ConfirmData

	sessionToken := c.Query("session")

	if err := c.ShouldBindJSON(&data); err != nil {
		errx.Handle(c, err)
		return
	}

	resp, err := h.authService.LoginConfirm(c.Request.Context(), &data, sessionToken, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		errx.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) RegistrationStart(c *gin.Context) {
	var data auth.AuthData

	if err := c.ShouldBindJSON(&data); err != nil {
		errx.Handle(c, errx.ErrInvalid)
		return
	}

	resp, err := h.authService.RegistrationStart(c.Request.Context(), &data, c.ClientIP())
	if err != nil {
		errx.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) RegistrationConfirm(c *gin.Context) {
	var data auth.ConfirmData

	sessionToken := c.Query("session")

	if err := c.ShouldBindJSON(&data); err != nil {
		errx.Handle(c, err)
		return
	}

	if err := h.authService.RegistrationConfirm(c.Request.Context(), &data, sessionToken, c.ClientIP()); err != nil {
		errx.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var data struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		errx.Handle(c, errx.ErrInvalid)
		return
	}

	token, err := h.tokenService.RefreshToken(c.Request.Context(), data.RefreshToken)
	if err != nil {
		errx.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *Handler) Logout(c *gin.Context) {
	accessToken := middleware.GetAccessToken(c)

	if err := h.tokenService.RevokeSession(c.Request.Context(), accessToken); err != nil {
		errx.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) LogoutAll(c *gin.Context) {
	accessToken := middleware.GetAccessToken(c)

	if err := h.tokenService.RevokeAllSession(c.Request.Context(), accessToken); err != nil {
		errx.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID := middleware.GetUserID(c)

	u, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		errx.Handle(c, err)
		return
	}

	c.JSON(http.StatusNoContent, u)
}

func (h *Handler) ResetPasswordStart(c *gin.Context) {
	var data auth.ResetPasswordStart

	if err := c.ShouldBindJSON(&data); err != nil {
		errx.Handle(c, errx.ErrInvalid)
		return
	}

	if err := h.authService.ResetPasswordStart(c.Request.Context(), &data, c.ClientIP()); err != nil {
		errx.Handle(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ResetPasswordConfirm(c *gin.Context) {
	var data auth.ResetPasswordConfirm

	sessionToken := c.Query("session")

	if err := c.ShouldBindJSON(&data); err != nil {
		errx.Handle(c, errx.ErrInvalid)
		return
	}

	if err := h.authService.ResetPasswordConfirm(c.Request.Context(), &data, sessionToken, c.ClientIP()); err != nil {
		errx.Handle(c, err)
		return
	}

	c.Status(http.StatusOK)
}
