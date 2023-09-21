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

// 倒计时插件
func Countdown(c *gin.Context) {
	cd := c.DefaultQuery("t", "2:00")
	c.HTML(http.StatusOK, "countdown.tmpl", gin.H{"cd": cd})
}
