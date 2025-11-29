package UserController

import (
	"fmt"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/controller"
	"github.com/StephenChristianW/go-movies-open/routers/Middlewares"
	"github.com/StephenChristianW/go-movies-open/security/IPManage"
	"github.com/StephenChristianW/go-movies-open/services"
	UserService "github.com/StephenChristianW/go-movies-open/services/User"
	"github.com/StephenChristianW/go-movies-open/utils/Jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
	Forget(c *gin.Context)
	ChangePassword(c *gin.Context)
	RenameUser(c *gin.Context)
}

var userService UserService.UserInterface = &UserService.UserSchema{}

func getUserInfo(c *gin.Context) *Jwt.UserClaims {
	return c.MustGet(Middlewares.UserInfo).(*Jwt.UserClaims)
}

type UserHandler struct {
}

// RegisterUser 注册新用户
// @Summary 注册新用户
// @Description 使用用户名、密码和邀请码注册用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body UserService.CreateUser true "用户注册信息"
// @Success 201 {object} swaggerResponse.SuccessResponseDoc "注册成功"
// @Failure 400 {object} userSwaggerResponse.RegisterUserErrorResponse "参数错误或注册失败"
// @Router /user/register [post]
func (UserHandler) RegisterUser(c *gin.Context) {
	var req UserService.CreateUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse("参数错误"))
		return
	}

	err = userService.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, controller.SuccessResponse(nil))
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码登录，返回不同状态码：
// @Description - code = -1：普通错误（返回累计错误次数）
// @Description - code = -2：达到锁定阈值（返回剩余解锁秒数，用于前端倒计时）
// @Description - code = -3：异常请求（例如请求过于频繁，疑似机器人）
// @Description - code = 1：登录成功（返回 token）
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body UserService.UserLogin true "登录请求参数"
// @Success 200 {object} userSwaggerResponse.LoginSuccessResponse "登录成功,返回token"
// @Failure 400 {object} userSwaggerResponse.LoginErrorResponse "登录失败，-1普通错误，-2限制登录倒计时，-3异常请求"
// @Router /user/login [post]
func (UserHandler) Login(c *gin.Context) {
	var req UserService.UserLogin
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse("参数错误"))
		return
	}
	data, err := userService.UserLogin(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}
	var failureResp UserService.LoginFailureResponse
	failureResp.FailedAttempts = data.FailedAttempts
	failureResp.LockSeconds = data.LockSeconds
	countMap := map[int]struct{}{3: {}, 5: {}, 7: {}, 10: {}}
	_, ok := countMap[failureResp.FailedAttempts]
	if !ok && failureResp.LockSeconds == 0 && data.FailedAttempts > 0 {
		c.JSON(http.StatusBadRequest, services.Response{
			Code: -1,
			Msg:  fmt.Sprintf("密码错误 %d 次", data.FailedAttempts),
			Data: failureResp,
		})
		return
	}
	if ok && failureResp.LockSeconds > 0 && data.FailedAttempts > 0 {

		c.JSON(http.StatusBadRequest, services.Response{
			Code: -2,
			Msg:  fmt.Sprintf("密码错误 %d 次, 请在 %d 秒后重试", data.FailedAttempts, data.LockSeconds),
			Data: failureResp,
		})
		return
	}
	if data.Token == "" {
		if data.FailedAttempts == 10 && !config.ProdEnv() {
			ip := c.ClientIP()
			IPManage.AddSuspiciousIp(req.Username, ip, "异常请求登录次数过多，涉嫌机器人刷密码")
		}
		failureResp.LockSeconds = 1
		c.JSON(http.StatusBadRequest, services.Response{
			Code: -3, // 自定义错误码，表示异常请求/机器人刷密码
			Msg:  "请求异常，请稍后重试或联系管理员",
			Data: nil,
		})
		return
	}
	var successResp UserService.LoginSuccessResponse
	successResp.Token = data.Token
	c.JSON(http.StatusOK, controller.SuccessResponse(successResp))
}

// Forget 忘记密码重置
// @Summary 忘记密码
// @Description 用户忘记密码时，通过用户名和邀请码重置密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body UserService.ForgetUser true "忘记密码请求体"
// @Success 200 {object} userSwaggerResponse.ForgetSuccessResponse "密码重置成功"
// @Failure 400 {object} userSwaggerResponse.ForgetErrorResponse "参数错误或重置失败"
// @Router /user/forget [put]
func (UserHandler) Forget(c *gin.Context) {
	var req UserService.ForgetUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse("参数错误"))
		return
	}
	err = userService.ForgetPassword(req.Username, req.Code, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(nil))
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 已登录用户通过旧密码修改为新密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body UserService.ChangePasswordSchema true "修改密码请求体"
// @Success 200 {object} userSwaggerResponse.ChangePasswordForgetSuccessResponse "密码修改成功"
// @Failure 400 {object} userSwaggerResponse.ChangePasswordErrorResponse "参数错误或修改失败"
// @Router /user/change_pwd [put]
// @Param Authorization header string true "用户 JWT Token"
func (UserHandler) ChangePassword(c *gin.Context) {
	var req UserService.ChangePasswordSchema
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse("参数错误"))
		return
	}
	userInfo := getUserInfo(c)
	err = userService.ChangePassword(userInfo.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(nil))
}

// RenameUser 修改当前登录用户的用户名
// @Summary 修改用户名
// @Description 修改当前登录用户的用户名
// @Tags 用户
// @Accept json
// @Produce json
// @Param username query string true "新的用户名"
// @Success 200 {object} userSwaggerResponse.RenameUserSuccessResponse "修改成功"
// @Failure 400 {object} userSwaggerResponse.RenameUserErrorResponse "参数错误 / 用户名已存在 / 用户不存在"
// @Router /user/rename [put]
// @Param Authorization header string true "用户 JWT Token"
func (UserHandler) RenameUser(c *gin.Context) {
	var req UserService.Rename
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse("参数错误"))
		return
	}
	userInfo := getUserInfo(c)
	err = userService.RenameUser(userInfo.UserID, req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse("修改成功"))
}
