package route

import (
	"DouYinService/controller"
	"DouYinService/douyin"
	"DouYinService/socket"
	"context"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Server struct {
	Bind      string
	DouYinUrl string
}

type config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Douyin struct {
		Room string `yaml:"room"`
	} `yaml:"douyin"`
}

var (
	r    = gin.Default()
	conf config
	ctx  = context.Background()
)

// 初始化配置
func initial_config() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}
}

// 初始化抖音协议
func initial_douyin(url string) {
	room, err := douyin.NewRoom(url)
	if err != nil {
		panic(err)
	}
	room.Connect()
}

// 初始化路由
func initial_route() {
	// 设置模版路径
	r.LoadHTMLGlob("template/**/*")
	// 设置静态文件
	r.Static("/static", "static/")
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
	r.GET("/plugin/countdown", controller.Countdown)
	// 插件管理
	r.GET("/manage/countdown", controller.CountdownManage)
	r.POST("/manage/countdown", controller.CountdownManageApi)
}

// 启动GIN服务器
func (s *Server) Start() {
	gin.SetMode(gin.ReleaseMode)
	// service.Initial(ctx)
	initial_config()
	initial_douyin(conf.Douyin.Room)
	initial_route()
	r.Run(":" + strconv.Itoa(conf.Server.Port))
}
