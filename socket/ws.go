package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	clients = make(map[*websocket.Conn]bool)
)

type WsMessage struct {
	// 消息类型 1 礼物，2 弹幕，3 点赞，4 进入房间
	Type int `json:"type"`
	// 礼物消息
	GiftName  string `json:"gift_name"`
	GiftCount uint64 `json:"gift_count"`
	GiftId    uint64 `json:"gift_id"`
	// 弹幕消息
	NickName   string `json:"nick_name"`
	HeadImgUrl string `json:"head_img"`
	Content    string `json:"content"`
	// 点赞数量
	Count uint64 `json:"count"`
}

func EndPoint(c *gin.Context) {
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	clients[conn] = true
}

func (m *WsMessage) Pust() {
	for client := range clients {
		err := client.WriteJSON(m)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
