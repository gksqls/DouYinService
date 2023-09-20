package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 木鱼插件
func Wooden(c *gin.Context) {
	c.HTML(http.StatusOK, "wooden.tmpl", nil)
}

// 弹幕插件
func Barrage(c *gin.Context) {
	c.HTML(http.StatusOK, "barrage.tmpl", nil)
}
