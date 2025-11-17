package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meszmate/zerra/internal/errx"
)

func (h *Handler) LoginStart(c *gin.Context) {

}

func (h *Handler) RefreshToken(c *gin.Context) {
	var data struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(data); err != nil {
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
