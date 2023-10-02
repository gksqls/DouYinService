package route

import (
	"DouYinService/controller"
	"DouYinService/douyin"
	"DouYinService/service"
	"DouYinService/socket"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Bind      string
	DouYinUrl string
}

var (
	r    = gin.Default()
	conf *service.TConfig
)

// 初始化配置
func initial_config() {
	service.Initial(context.Background())
	conf = new(service.TConfig)
	_, err := service.Db.Where("id=?", 1).Get(conf)
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
		// 加班倒计时管理
		manage.GET("/overtime", controller.OvertimeManage)
		manage.POST("/overtime", controller.OvertimeManageApi)
		manage.POST("/overtime/save", controller.OvertimeManageSaveApi)
		manage.POST("/overtime/gift", controller.OvertimeGiftManageApi)
		// 配置管理
		manage.POST("/config", controller.ConfigManageApi)
		manage.POST("/config/update", controller.ConfigManageUpdateApi)
	}
}

// 启动GIN服务器
func (s *Server) Start() {
	gin.SetMode(gin.ReleaseMode)
	initial_config()
	initial_douyin(conf.DouyinRoom)
	initial_route()
	r.Run(":" + strconv.Itoa(conf.ServerPort))
}
