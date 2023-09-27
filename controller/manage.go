package controller

import (
	"DouYinService/controller/result"
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
	c.JSON(http.StatusOK, result.Success())
}

// 礼物列表管理页面
func GiftListManage(c *gin.Context) {
	c.HTML(http.StatusOK, "gift_manage.tmpl", nil)
}

// 礼物列表管理API
func GiftListManageApi(c *gin.Context) {
	c.JSON(http.StatusOK, result.Success())
}

// 加班倒计时管理
func OvertimeManage(c *gin.Context) {
	c.HTML(http.StatusOK, "overtime_manage.tmpl", nil)
}

// 加班倒计时管理
func OvertimeManageApi(c *gin.Context) {
	overtime := &service.Overtime{Context: c}
	overtime.Set()
	c.JSON(http.StatusOK, result.Success())
}
