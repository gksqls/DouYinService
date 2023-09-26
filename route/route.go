package route

import (
	"DouYinService/controller"
	"DouYinService/douyin"
	"DouYinService/socket"
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
	// 礼物模版
	gift := r.Group("/gift")
	{
		// 跳舞礼物模版
		gift.GET("/left", controller.GiftDanceLeft)
		gift.GET("/foot", controller.GiftDanceFoot)
		// 游戏礼物模版
		gift.GET("/game/left", controller.GiftGameLeft)
		gift.GET("/game/foot", controller.GiftGameFoot)
	}
	// 插件
	plugin := r.Group("/plugin")
	{
		plugin.GET("/wooden", controller.Wooden)
		plugin.GET("/barrage", controller.Barrage)
		plugin.GET("/countdown", controller.Countdown)
		plugin.GET("/overtime", controller.Overtime)
	}
	// 插件管理
	manage := r.Group("/manage")
	{
		// 倒计时插件管理
		manage.GET("/countdown", controller.CountdownManage)
		manage.POST("/countdown", controller.CountdownManageApi)
		// 送礼列表管理
		manage.GET("/gift", controller.GiftListManage)
		manage.POST("/gift", controller.GiftListManageApi)
	}
}

// 启动GIN服务器
func (s *Server) Start() {
	gin.SetMode(gin.ReleaseMode)
	// service.Initial(context.Background())
	initial_config()
	initial_douyin(conf.Douyin.Room)
	initial_route()
	r.Run(":" + strconv.Itoa(conf.Server.Port))
}
