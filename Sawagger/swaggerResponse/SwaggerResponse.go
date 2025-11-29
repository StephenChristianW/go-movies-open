package swaggerResponse

// SuccessResponseDoc 用于 swaggerResponse 成功返回
type SuccessResponseDoc struct {
	Code int    `json:"code" example:"1"`   // 状态码
	Msg  string `json:"msg" example:"ok"`   // 信息
	Data string `json:"data" example:"any"` // 数据
}

// ErrorResponseDoc 用于 swaggerResponse 失败返回
type ErrorResponseDoc struct {
	Code int    `json:"code" example:"-1"`   // 状态码
	Msg  string `json:"msg"  example:"参数错误"` // 信息
	Data string `json:"data" example:"null"` // 数据
}
