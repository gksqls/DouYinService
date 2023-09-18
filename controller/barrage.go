package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Barrage(c *gin.Context) {
	c.HTML(http.StatusOK, "barrage.tmpl", nil)
}
