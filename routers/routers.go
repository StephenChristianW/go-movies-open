package routers

import (
	"github.com/StephenChristianW/go-movies-open/controller/AdminController"
	"github.com/StephenChristianW/go-movies-open/controller/InviteController"
	"github.com/StephenChristianW/go-movies-open/controller/UserController"
	"github.com/StephenChristianW/go-movies-open/routers/Middlewares"
	"github.com/gin-gonic/gin"
)

// ==================================== 路由注册 ====================================

// registerInviteRoutes 注册 InviteCode 模块的路由
// 参数 rg: *gin.RouterGroup, 路由分组对象
func registerInviteRoutes(rg *gin.RouterGroup) {
	var handler InviteController.InviteCodeController = &InviteController.InviteCodeHandler{}

	rg.POST("/create_one", handler.GenerateAndInsertCode)    // 生成单个邀请码
	rg.POST("/create_batch", handler.GenerateAndInsertCodes) // 批量生成邀请码
	rg.POST("/reback", handler.ReBackInviteCodes)            // 根据ID回查邀请码
	rg.POST("/list", handler.ListInviteCodes)                // 分页获取邀请码列表
	rg.DELETE("/delete", handler.DeleteInviteCode)           // 删除邀请码
	rg.PUT("/update", handler.UpdateInviteCode)              // 更新邀请码状态
	rg.GET("/get", handler.GetInviteCode)                    // GET获取邀请码
}

// registerAdminRoutes 注册 Admin 模块的路由
// 参数 rg: *gin.RouterGroup, 路由分组对象
func registerAdminRoutes(rg *gin.RouterGroup) {
	var handler AdminController.AdminController = &AdminController.AdminHandler{}
	rg.POST("/rlogin", handler.Remember)  // 管理员记住密码登录
	rg.POST("/login", handler.AdminLogin) // 管理员登录
	rg.POST("/logout", handler.LogOut)    // 管理员登出

	rg.POST("/list", handler.AdminList)                // 获取管理员列表
	rg.POST("/create", handler.CreateAdmin)            // 创建管理员
	rg.PUT("/ban", handler.BannedAdminUser)            // 禁用管理员
	rg.PUT("/active", handler.ActiveAdminUser)         // 启用管理员
	rg.PUT("/change_pwd", handler.ChangeAdminPassword) // 修改管理员密码
	rg.DELETE("/delete", handler.DeleteAdmin)          // 删除管理员
}

func registerUserRoutes(rg *gin.RouterGroup) {
	var handler UserController.UserController = &UserController.UserHandler{}
	rg.POST("/register", handler.RegisterUser)    // 创建用户
	rg.POST("/login", handler.Login)              // 创建用户
	rg.PUT("/forget", handler.Forget)             // 忘记密码
	rg.PUT("/change_pwd", handler.ChangePassword) // 创建用户
	rg.PUT("/rename", handler.RenameUser)         // 创建用户
}

// registerRoutes 注册所有模块路由
// - 创建路由分组（/invite, /admin）
// - 挂载鉴权中间件（AdminAuthMiddleware）
// - 调用各模块路由注册函数
func registerRoutes() {
	// 路由分组
	inviteGroup := router.Group("/invite")
	adminGroup := router.Group("/admin")
	userGroup := router.Group("/user")
	// 挂载中间件
	adminGroup.Use(Middlewares.AdminAuthMiddleware())
	inviteGroup.Use(Middlewares.AdminAuthMiddleware())
	userGroup.Use(Middlewares.UserAuthMiddleware())

	// 注册模块路由
	registerInviteRoutes(inviteGroup)
	registerAdminRoutes(adminGroup)
	registerUserRoutes(userGroup)
}
