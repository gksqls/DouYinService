package service

import (
	"DouYinService/socket"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Overtime struct {
	Context *gin.Context
}

func (o *Overtime) Set() {
	jsonMap := make(map[string]interface{})
	o.Context.BindJSON(&jsonMap)
	t, _ := strconv.Atoi(jsonMap["type"].(string))
	cd, _ := strconv.Atoi(jsonMap["time"].(string))
	message := &socket.WsMessage{Type: t, Countdown: cd}
	message.Pust()
}
