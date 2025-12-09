package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/meszmate/zerra/internal/errx"
)

func (h *Handler) UpdateAvatar(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		errx.Handle(c, errx.ErrImage)
		return
	}

	fmt.Println(f.Size)
	return
}
