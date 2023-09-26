package controller

import (
	"DouYinService/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	index := &service.Index{Context: c}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"PList": index.LoadPluginList(),
	})
}
