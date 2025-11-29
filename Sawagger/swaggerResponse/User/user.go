package userSwaggerResponse

type RegisterUserSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"nil"`
}
type RegisterUserErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}

type LoginSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"{token: eyJhbGciOiJIUzI1N...}"`
}
type LoginErrorResponse struct {
	Code string `json:"code" example:"-1/-2/-3"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"-1:null;-2:{num: 3, seconds: 60}; -3:null"`
}

type ForgetSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"null"`
}
type ForgetErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}

type ChangePasswordForgetSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"null"`
}
type ChangePasswordErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}
type RenameUserSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"修改成功"`
}
type RenameUserErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"参数错误 / 用户名已存在 / 用户不存在"`
	Data string `json:"data" example:"null"`
}
