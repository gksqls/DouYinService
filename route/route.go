package route

import (
	"DouYinService/controller"
	"DouYinService/douyin"
	"DouYinService/socket"

	"github.com/gin-gonic/gin"
)

var (
	r = gin.Default()
)

type Server struct {
	Bind      string
	DouYinUrl string
}

func initial_douyin(url string) {
	room, err := douyin.NewRoom(url)
	if err != nil {
		panic(err)
	}
	room.Connect()
}

func initial_route() {
	// 主页
	r.GET("/", controller.Index)
	r.GET("/index", controller.Index)
	// websocket endPoint
	r.GET("/nscd", socket.EndPoint)
	// 跳舞礼物模版
	r.GET("/gift/left", controller.GiftDanceLeft)
	r.GET("/gift/foot", controller.GiftDanceFoot)
	// 游戏礼物模版
	r.GET("/gift/game/left", controller.GiftGameLeft)
	r.GET("/gift/game/foot", controller.GiftGameFoot)
	// 插件
	r.GET("/plugin/wooden", controller.Wooden)
	r.GET("/plugin/barrage", controller.Barrage)
}

func (s *Server) Initial() {
	gin.SetMode(gin.ReleaseMode)
	r.LoadHTMLGlob("template/**/*")
	r.Static("/static", "static/")
	initial_douyin(s.DouYinUrl)
	initial_route()
}

func (s *Server) Start() {
	s.Initial()
	r.Run(s.Bind)
}
