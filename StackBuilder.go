package main

import (
	_ "github.com/StephenChristianW/go-movies-open/Sawagger/docs" // 引入 swaggerResponse 自动生成的文档
	_ "github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/config/initRootAdmin"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	"github.com/StephenChristianW/go-movies-open/routers"
	InviteCodeService "github.com/StephenChristianW/go-movies-open/services/System/InviteCode"
	"github.com/StephenChristianW/go-movies-open/task"
)

func main() {
	// ================== 初始化 MongoDB 连接 ==================
	// 获取全局 MongoDB 客户端并确保程序退出时自动关闭连接
	defer db.MongoClient()()

	// 初始化根管理员
	initRootAdmin.InitRootAdmins()

	// ================== 定时任务 ==================
	// 每日零点执行一次 CheckExpiredCodes，处理过期邀请码
	task.RunDaily("CheckExpiredCodes", InviteCodeService.CheckExpiredCodes)

	// ================== 启动 HTTP 服务 ==================
	// 启动 Gin HTTP 服务器，包括路由注册、swaggerResponse 和路由打印
	routers.RunServer()

}
