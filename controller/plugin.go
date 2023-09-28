package controller

import (
	"DouYinService/service"
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
	c.HTML(http.StatusOK, "countdown.tmpl", nil)
}

// 加班插件
func Overtime(c *gin.Context) {
	overtime := &service.Overtime{Context: c}
	c.HTML(http.StatusOK, "overtime.tmpl", gin.H{"overtime": overtime.Find()})
}
