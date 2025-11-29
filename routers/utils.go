package routers

import (
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
)

const (
	routesLog = "[ROUTES] "
	serverLog = "[SERVER] "
)

func runServer() {
	// 启动服务
	port := config.Addr
	if err := router.Run(port); err != nil {
		log.Fatalf(serverLog+"启动失败: %v\n", err)
	}
}

// setSwagger 配置 swaggerResponse 路由，仅在非生产环境
func setSwagger() {
	if config.ProdEnv() {
		return
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println(serverLog + "Server running at http://localhost:" + config.Port)
	log.Println(serverLog + "swaggerResponse UI: http://localhost:" + config.Port + "/swagger/index.html")
}

// printRoutes 打印已注册路由信息
func printRoutes() {
	log.Println()
	log.Println(routesLog + "已注册路由列表:")
	for _, ri := range router.Routes() {
		log.Printf(" - %s %s -> %s\n", ri.Method, ri.Path, ri.Handler)
	}
}

// setGinMode 根据环境设置 Gin 运行模式
func setGinMode() {
	mode := gin.DebugMode
	envDesc := "本地环境"
	if config.ProdEnv() {
		mode = gin.ReleaseMode
		envDesc = "生产环境"
		printRoutes() // 打印已注册路由信息
	}
	gin.SetMode(mode)
	log.Printf("[ENV] %s, Gin 设置为 %s\n", envDesc, mode)
	log.Printf("[ENV] 当前环境: %s\n", config.Env)
}
