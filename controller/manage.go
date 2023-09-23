package controller

import (
	"DouYinService/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 倒计时管理页面
func CountdownManage(c *gin.Context) {
	c.HTML(http.StatusOK, "countdown_manage.tmpl", nil)
}

// 倒计时管理API
func CountdownManageApi(c *gin.Context) {
	countdown := &service.Countdown{Context: c}
	countdown.Start()
	c.JSON(http.StatusOK, success())
}
