package adminSwaggerResponse

type AdminRememberSuccessResponse struct {
	Code    int    `json:"code" example:"1"`
	Message string `json:"message" example:"ok"`
	Data    string `json:"data" example:"{token:eyJhbGciOiJIUzI1N...,rt:eyJhbGciOiJIUzI1N...}"`
}
type AdminRememberErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"err_msg"`
	Data string `json:"data" example:"null"`
}

type CreateAdminSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"管理员 xxx 已创建"`
}

type CreateAdminErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"管理员 xxx 已存在"`
	Data string `json:"data" example:"null"`
}

type AdminLoginSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"eyJhbGciOiJIUzI1N..."`
}
type AdminLoginErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"密码错误"`
	Data string `json:"data" example:"null"`
}

type BannedAdminUserSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"管理员 xxx 已被禁用"`
}

type BannedAdminUserErrorResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}
type ActiveAdminUserSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"管理员 xxx 已启用"`
}

type ActiveAdminUserErrorResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}

type LogOutSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"登出成功"`
}
type LogOutErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}

type AdminListSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"{page:1, page_size: 10, total_pages: 100,data: [id1_Detail, id2_Detail, id3_Detail]}"`
}
type AdminListErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}

type DeleteAdminSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"管理员 xxx 已删除"`
}
type DeleteAdminErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}

type ChangeAdminPasswordSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data" example:"管理员 xxx 密码已更新"`
}
type ChangeAdminPasswordErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"err msg"`
	Data string `json:"data" example:"null"`
}
