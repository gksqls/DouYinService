package service

import (
	"DouYinService/socket"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Countdown struct {
	Context *gin.Context
}

func (c *Countdown) Start() {
	jsonMap := make(map[string]interface{})
	c.Context.BindJSON(&jsonMap)
	cd, _ := strconv.Atoi(jsonMap["t"].(string))
	message := &socket.WsMessage{
		Type:      101,
		Countdown: cd,
	}
	message.Pust()
}
