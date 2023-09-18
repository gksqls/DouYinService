package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Wooden(c *gin.Context) {
	c.HTML(http.StatusOK, "wooden.tmpl", nil)
}
