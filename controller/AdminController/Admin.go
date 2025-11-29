package AdminController

import (
	"fmt"
	"github.com/StephenChristianW/go-movies-open/controller"
	"github.com/StephenChristianW/go-movies-open/routers/Middlewares"
	AdminService2 "github.com/StephenChristianW/go-movies-open/services/System/Admin"
	"github.com/StephenChristianW/go-movies-open/utils/Jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminController 管理员接口定义
type AdminController interface {
	CreateAdmin(c *gin.Context)
	AdminLogin(c *gin.Context)
	BannedAdminUser(c *gin.Context)
	LogOut(c *gin.Context)
	AdminList(c *gin.Context)
	DeleteAdmin(c *gin.Context)
	ChangeAdminPassword(c *gin.Context)
	Remember(c *gin.Context)
	ActiveAdminUser(c *gin.Context)
}

// AdminHandler 实现 AdminController
type AdminHandler struct{}

var adminService AdminService2.AdminInterface = &AdminService2.AdminService{}

// -------------------- 工具函数 --------------------

// bindAndValidate 绑定请求参数并校验，如果失败直接返回 400
func bindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBind(obj); err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse("参数错误"))
		return false
	}
	return true
}

// getAdminInfo 从 Context 获取当前管理员信息
func getAdminInfo(c *gin.Context) *Jwt.AdminClaims {
	return c.MustGet(Middlewares.AdminInfo).(*Jwt.AdminClaims)
}

// -------------------- 控制器方法 --------------------

// CreateAdmin 创建管理员
// @Summary 创建管理员
// @Description 创建新的管理员账户
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param body body AdminService.AdminCreate true "用户名、密码、角色"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} adminSwaggerResponse.CreateAdminSuccessResponse
// @Failure 422 {object} adminSwaggerResponse.CreateAdminErrorResponse
// @Router /admin/create [post]
// 说明：
//   - 功能：接收用户名、密码和角色ID，调用 service 层创建管理员
//   - 输入：JSON body 包含 Username, Password, Role
//   - 输出：成功返回 200 + 成功提示，失败返回 422 + 错误信息
func (*AdminHandler) CreateAdmin(c *gin.Context) {
	var req AdminService2.AdminCreate
	if !bindAndValidate(c, &req) {
		return
	}
	if err := adminService.CreateAdmin(req.Username, req.Password, req.Role); err != nil {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(fmt.Sprintf("管理员 %s 已创建", req.Username)))
}

// Remember 管理员登录并记住会话
// @Summary 管理员记住登录
// @Description 登录后返回访问令牌和刷新令牌
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param body body AdminService.LoginCredentialsRequest true "登录凭证"
// @Success 200 {object} adminSwaggerResponse.AdminRememberSuccessResponse
// @Failure 401 {object} adminSwaggerResponse.AdminRememberErrorResponse
// @Router /admin/remember [post]
// 说明：
//   - 功能：支持用户名/密码或刷新令牌登录，返回短期访问令牌和长期刷新令牌
//   - 输入：JSON body 包含 Username, Password, DeviceId, RT, AuthStatus, IP
//   - 输出：成功返回 200 + {Token, RToken}，失败返回 401 + 错误信息
func (*AdminHandler) Remember(c *gin.Context) {
	var req AdminService2.LoginCredentialsRequest
	if !bindAndValidate(c, &req) {
		return
	}
	req.IP = c.ClientIP()
	data, err := adminService.LoginAndRemember(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(data))
}

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 使用用户名密码登录，返回 JWT
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param body body AdminService.AdminLoginRequest true "用户名和密码"
// @Success 200 {object} adminSwaggerResponse.AdminLoginSuccessResponse
// @Failure 401 {object} adminSwaggerResponse.AdminLoginErrorResponse
// @Router /admin/login [post]
// 说明：
//   - 功能：使用用户名密码进行登录，生成短期有效的 JWT
//   - 输入：JSON body 包含 Username, Password, DeviceId, Remember, IP
//   - 输出：成功返回 200 + Token，失败返回 401 + 错误信息
func (*AdminHandler) AdminLogin(c *gin.Context) {
	var req AdminService2.AdminLoginRequest
	if !bindAndValidate(c, &req) {
		return
	}
	token, err := adminService.AdminLogin(req.Username, req.Password, c.ClientIP(), req.DeviceId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(token))
}

// LogOut 管理员登出
// @Summary 管理员登出
// @Description 登出当前管理员
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} adminSwaggerResponse.LogOutSuccessResponse
// @Failure 401 {object} adminSwaggerResponse.LogOutErrorResponse
// @Router /admin/logout [post]
// 说明：
//   - 功能：删除当前管理员的 JWT token，实现登出
//   - 输入：无
//   - 输出：成功返回 200 + 成功提示，失败返回 401 + 错误信息
func (*AdminHandler) LogOut(c *gin.Context) {
	adminInfo := getAdminInfo(c)
	if err := adminService.LogOut(adminInfo.AdminID); err != nil {
		c.JSON(http.StatusUnauthorized, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse("登出成功"))
}

// ActiveAdminUser 启用管理员
// @Summary 启用管理员
// @Description 启用指定管理员（恢复被禁用的账户）
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param body body AdminService.BannedUserRequest true "被启用用户名"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} adminSwaggerResponse.ActiveAdminUserSuccessResponse
// @Failure 422 {object} adminSwaggerResponse.ActiveAdminUserErrorResponse
// @Router /admin/active [put]
//
// 说明：
//   - 功能：将指定管理员状态置为启用
//   - 输入：JSON body 包含 BannedUsername
//   - 输出：成功返回 200 + 成功提示，失败返回 422 + 错误信息
func (*AdminHandler) ActiveAdminUser(c *gin.Context) {
	var req AdminService2.BannedUserRequest
	if !bindAndValidate(c, &req) {
		return
	}
	if err := adminService.ActiveAdminUser(req.BannedUsername); err != nil {
		c.JSON(http.StatusUnauthorized, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(fmt.Sprintf("管理员 %s 已启用", req.BannedUsername)))
}

// BannedAdminUser 禁用管理员
// @Summary 禁用管理员
// @Description 禁用指定管理员
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param body body AdminService.BannedUserRequest true "被禁用用户名"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} adminSwaggerResponse.BannedAdminUserSuccessResponse
// @Failure 422 {object} adminSwaggerResponse.BannedAdminUserErrorResponse
// @Router /admin/ban [put]
// 说明：
//   - 功能：将指定管理员状态置为禁用，并删除 Redis token
//   - 输入：JSON body 包含 BannedUsername
//   - 输出：成功返回 200 + 成功提示，失败返回 422 + 错误信息
func (*AdminHandler) BannedAdminUser(c *gin.Context) {
	var req AdminService2.BannedUserRequest
	if !bindAndValidate(c, &req) {
		return
	}
	if err := adminService.BannedAdminUser(req.BannedUsername); err != nil {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(fmt.Sprintf("管理员 %s 已被禁用", req.BannedUsername)))
}

// AdminList 获取管理员列表
// @Summary 获取管理员列表
// @Description 获取分页管理员列表，可过滤用户名和状态
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param query query AdminService.AdminList true "分页及过滤参数"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} adminSwaggerResponse.AdminListSuccessResponse
// @Failure 422 {object} swaggerResponse.ErrorResponseDoc
// @Failure 500 {object} adminSwaggerResponse.AdminListErrorResponse
// @Router /admin/list [post]
// 说明：
//   - 功能：分页获取管理员列表，可通过用户名、状态、角色过滤
//   - 输入：JSON query 包含 Page, PageSize, Username, Status, Deleted, Role, IDs
//   - 输出：成功返回 200 + 分页数据，失败返回 422 或 500 + 错误信息
func (*AdminHandler) AdminList(c *gin.Context) {
	var req AdminService2.AdminList
	if !bindAndValidate(c, &req) {
		return
	}
	req.Page, req.PageSize = controller.PageSet(req.Page, req.PageSize)
	adminList, err := adminService.AdminList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(adminList))
}

// DeleteAdmin 删除管理员
// @Summary 删除管理员
// @Description 根据用户名删除管理员
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param body body AdminService.DeletedUserRequest true "被删除用户名"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} adminSwaggerResponse.DeleteAdminSuccessResponse
// @Failure 422 {object} adminSwaggerResponse.DeleteAdminErrorResponse
// @Router /admin/delete [delete]
// 说明：
//   - 功能：逻辑删除指定管理员账号，并删除 Redis token
//   - 输入：JSON body 包含 DelUsername
//   - 输出：成功返回 200 + 成功提示，失败返回 422 + 错误信息
func (*AdminHandler) DeleteAdmin(c *gin.Context) {
	var req AdminService2.DeletedUserRequest
	if !bindAndValidate(c, &req) {
		return
	}
	adminInfo := getAdminInfo(c)
	if err := adminService.DeleteAdmin(adminInfo.AdminUsername, req.DelUsername); err != nil {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(fmt.Sprintf("管理员 %s 已删除", req.DelUsername)))
}

// ChangeAdminPassword 修改管理员密码
// @Summary 修改管理员密码
// @Description 修改指定管理员密码
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param body body AdminService.ChangeAdminPasswordRequest true "用户名和新密码"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} adminSwaggerResponse.ChangeAdminPasswordSuccessResponse
// @Failure 422 {object} adminSwaggerResponse.ChangeAdminPasswordErrorResponse
// @Router /admin/change_pwd [put]
// 说明：
//   - 功能：更新指定管理员的密码，并删除其 Redis token
//   - 输入：JSON body 包含 Username, NewPassword
//   - 输出：成功返回 200 + 成功提示，失败返回 422 + 错误信息
func (*AdminHandler) ChangeAdminPassword(c *gin.Context) {
	var req AdminService2.ChangeAdminPasswordRequest
	if !bindAndValidate(c, &req) {
		return
	}
	if err := adminService.ChangeAdminPassword(req.Username, req.NewPassword); err != nil {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(fmt.Sprintf("管理员 %s 密码已更新", req.Username)))
}
