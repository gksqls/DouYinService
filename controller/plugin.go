package controller

import (
	"net/http"
	"strconv"
	"time"

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
	t, err := strconv.ParseInt(c.DefaultQuery("t", "120"), 10, 64)
	if err != nil {
		c.HTML(http.StatusOK, "countdown.tmpl", gin.H{"cd": "02:00"})
	} else {
		c.HTML(http.StatusOK, "countdown.tmpl", gin.H{"cd": time.Unix(t, 0).Format("04:05")})
	}
}
