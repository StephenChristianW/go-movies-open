package routers

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine // 全局 Gin 引擎实例，可供其他地方引用

// RunServer 启动 HTTP 服务
// 1. 根据当前环境设置 Gin 模式（Debug/Release）
// 2. 初始化 Gin 引擎
// 3. 注册路由分组及模块路由
// 4. 配置 swaggerResponse（仅非生产环境）
// 5. 打印已注册路由信息
// 6. 启动 HTTP 服务监听端口
func RunServer() {
	setGinMode()           // 设置 Gin 模式
	router = gin.Default() // 初始化 Gin 引擎
	registerRoutes()       // 注册路由分组及模块路由
	setSwagger()           // 配置 swaggerResponse，仅非生产环境
	htmlController()
	setIcon()
	runServer() // 启动 HTTP 服务
}
