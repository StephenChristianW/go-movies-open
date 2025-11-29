package InviteController

import (
	"errors"
	"github.com/StephenChristianW/go-movies-open/controller"
	InviteCodeService2 "github.com/StephenChristianW/go-movies-open/services/System/InviteCode"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// InviteCodeController 邀请码控制器接口
type InviteCodeController interface {
	GenerateAndInsertCode(c *gin.Context)
	GenerateAndInsertCodes(c *gin.Context)
	ReBackInviteCodes(c *gin.Context)
	ListInviteCodes(c *gin.Context)
	DeleteInviteCode(c *gin.Context)
	UpdateInviteCode(c *gin.Context)
	GetInviteCode(c *gin.Context)
}

// InviteCodeHandler 实现
type InviteCodeHandler struct{}

var inviteCodeService InviteCodeService2.InviteCodeInterface = &InviteCodeService2.CreateInviteCode{}

// -------------------- 请求体定义 --------------------

type GenerateCodesRequest struct {
	Num int `json:"num" binding:"required"`
}

type InviteCodesRequest struct {
	CodeIds []string `json:"codeIds" binding:"required"`
}

// -------------------- 控制器方法 --------------------

// GenerateAndInsertCode 生成单个邀请码
// @Summary 生成单个邀请码
// @Description 生成一个新的邀请码
// @Tags 邀请码管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} inviteCodeSwaggerResponse.GenerateAndInsertCodeSuccessResponse
// @Failure 500 {object} inviteCodeSwaggerResponse.GenerateAndInsertCodeErrorResponse
// @Router /invite/create_one [post]
func (InviteCodeHandler) GenerateAndInsertCode(c *gin.Context) {
	data, err := inviteCodeService.GenerateAndInsertCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(data))
}

// GenerateAndInsertCodes 2批量生成邀请码
// @Summary 批量生成邀请码
// @Description 批量生成指定数量的邀请码
// @Tags 邀请码管理
// @Accept json
// @Produce json
// @Param body body GenerateCodesRequest true "生成数量"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} inviteCodeSwaggerResponse.GenerateAndInsertCodesSuccessResponse
// @Failure 422 {object} swaggerResponse.ErrorResponseDoc
// @Failure 500 {object} inviteCodeSwaggerResponse.GenerateAndInsertCodesErrorResponse
// @Router /invite/create_batch [post]
func (InviteCodeHandler) GenerateAndInsertCodes(c *gin.Context) {
	var req GenerateCodesRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Num <= 0 {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse("参数错误"))
		return
	}
	data, err := inviteCodeService.GenerateAndInsertCodes(req.Num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(data))
}

// ReBackInviteCodes 批量回查邀请码
// @Summary 批量回查邀请码
// @Description 根据 ID 查询邀请码信息
// @Tags 邀请码管理
// @Accept json
// @Produce json
// @Param body body InviteCodesRequest true "邀请码 ID 列表"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} swaggerResponse.SuccessResponseDoc
// @Failure 422 {object} swaggerResponse.ErrorResponseDoc
// @Failure 400 {object} swaggerResponse.ErrorResponseDoc
// @Router /invite/reback [post]
func (InviteCodeHandler) ReBackInviteCodes(c *gin.Context) {
	var req InviteCodesRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.CodeIds) == 0 {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse("参数错误"))
		return
	}
	data, err := inviteCodeService.ReBackInviteCodes(req.CodeIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(data))
}

// ListInviteCodes 分页获取邀请码列表
// @Summary 分页获取邀请码列表
// @Description 根据分页参数和状态获取邀请码列表
// @Tags 邀请码管理
// @Accept json
// @Produce json
// @Param body body InviteCodeService.CreateInviteFilter true "分页参数和过滤条件"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} inviteCodeSwaggerResponse.ListInviteCodesSuccessResponse
// @Failure 422 {object} swaggerResponse.ErrorResponseDoc
// @Failure 500 {object} inviteCodeSwaggerResponse.ListInviteCodesErrorResponse
// @Router /invite/list [post]
func (InviteCodeHandler) ListInviteCodes(c *gin.Context) {
	var req InviteCodeService2.CreateInviteFilter
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse("参数错误"))
		return
	}
	req.Page, req.PageSize = controller.PageSet(req.Page, req.PageSize)
	data, err := inviteCodeService.ListInviteCodes(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(data))
}

// DeleteInviteCode 删除邀请码
// @Summary 删除邀请码
// @Description 根据 ID 列表删除邀请码
// @Tags 邀请码管理
// @Accept json
// @Produce json
// @Param body body InviteCodeService.CreateInviteFilter true "ID 列表"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} inviteCodeSwaggerResponse.DeleteInviteCodeSuccessResponse
// @Failure 422 {object} swaggerResponse.ErrorResponseDoc
// @Failure 400 {object} inviteCodeSwaggerResponse.DeleteInviteCodeErrorResponse
// @Router /invite/delete [delete]
func (InviteCodeHandler) DeleteInviteCode(c *gin.Context) {
	var req InviteCodeService2.CreateInviteFilter
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse("未提供有效的 ID 列表"))
		return
	}
	data, err := inviteCodeService.DeleteInviteCode(req.IDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(data))
}

// UpdateInviteCode 更新邀请码状态
// @Summary 更新邀请码状态
// @Description 更新指定邀请码的状态
// @Tags 邀请码管理
// @Accept json
// @Produce json
// @Param body body InviteCodeService.UpdateInviteCode true "更新数据"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} inviteCodeSwaggerResponse.UpdateInviteCodeSuccessResponse
// @Failure 422 {object} swaggerResponse.ErrorResponseDoc
// @Failure 400 {object} inviteCodeSwaggerResponse.UpdateInviteCodeErrorResponse
// @Router /invite/update [put]
func (InviteCodeHandler) UpdateInviteCode(c *gin.Context) {
	var req InviteCodeService2.UpdateInviteCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse("参数错误"))
		return
	}
	data, err := inviteCodeService.UpdateInviteCodes(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, controller.SuccessResponse(data))
}

// GetInviteCode 根据 codeID 或 code 获取邀请码
// @Summary 获取邀请码
// @Description 根据 ObjectID 或 code 字段获取邀请码信息
// @Tags 邀请码管理
// @Accept json
// @Produce json
// @Param code query string true "邀请码"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} inviteCodeSwaggerResponse.GetInviteCodeErrorResponse "成功返回邀请码信息"
// @Failure 400 {object} inviteCodeSwaggerResponse.GetInviteCodeSuccessResponse "查询失败或未找到"
// @Failure 422 {object} swaggerResponse.ErrorResponseDoc "参数错误"
// @Router /invite/get [get]
func (InviteCodeHandler) GetInviteCode(c *gin.Context) {
	var req InviteCodeService2.InviteCodeRequest
	if err := c.ShouldBindQuery(&req); err != nil || req.Code == "" {
		c.JSON(http.StatusUnprocessableEntity, controller.ErrorResponse("参数错误"))
		return
	}

	data, err := inviteCodeService.GetInviteCode(req.Code)
	if err != nil {
		var errStr string
		if len(req.Code) == 24 {
			errStr = " id: " + req.Code
		} else if len(req.Code) == 8 {
			errStr = ": " + req.Code
		} else {
			errStr = "-验证码长度错误: " + req.Code
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusBadRequest, controller.ErrorResponse("未找到验证码"+errStr))
			return
		}
		c.JSON(http.StatusBadRequest, controller.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, controller.SuccessResponse(data))
}
