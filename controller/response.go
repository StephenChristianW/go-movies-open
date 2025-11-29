package controller

import (
	"github.com/StephenChristianW/go-movies-open/services"
)

//  -------------------- 通用响应封装 --------------------

func ErrorResponse(msg string) services.Response {
	return services.Response{
		Code: -1,
		Msg:  msg,
		Data: nil,
	}
}

func SuccessResponse(data interface{}) services.Response {
	return services.Response{
		Code: 1,
		Msg:  "ok",
		Data: data,
	}
}

const RequestParameterIncorrect = "RequestParameterIncorrect" // 请求参数不正确
