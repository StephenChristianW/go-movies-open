package inviteCodeSwaggerResponse

type GenerateAndInsertCodeSuccessResponse struct {
	Code int    `json:"code" example:"1"`         // 状态码
	Msg  string `json:"msg"  example:"ok"`        // 信息
	Data string `json:"data" example:"id string"` // 数据
}
type GenerateAndInsertCodeErrorResponse struct {
	Code int    `json:"code" example:"-1"`        // 状态码
	Msg  string `json:"msg"  example:"error msg"` // 信息
	Data string `json:"data" example:"null"`      // 数据
}

type GenerateAndInsertCodesSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data int    `json:"data"  example:"10"`
}
type GenerateAndInsertCodesErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"error msg"`
	Data string `json:"data" example:"null"`
}
type ReBackInviteCodesSuccessResponse struct {
	Code int      `json:"code" example:"1"`
	Msg  string   `json:"msg"  example:"ok"`
	Data []string `json:"data"  example:"[id1,id2,id3]"`
}
type ReBackInviteCodesErrorResponse struct {
	Code int      `json:"code" example:"-1"`
	Msg  string   `json:"msg"  example:"error msg"`
	Data []string `json:"data" example:"[]"`
}
type ListInviteCodesSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data"  example:"{page:1, pagesize: 10, totalPages: 100,data: [id1_Detail, id2_Detail, id3_Detail]}"`
}

type ListInviteCodesErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"error msg"`
	Data string `json:"data" example:"null"`
}
type DeleteInviteCodeSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data int    `json:"data"  example:"10"`
}
type DeleteInviteCodeErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"error msg"`
	Data string `json:"data" example:"null"`
}

type UpdateInviteCodeSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data int    `json:"data"  example:"10"`
}
type UpdateInviteCodeErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"error msg"`
	Data string `json:"data" example:"null"`
}
type GetInviteCodeSuccessResponse struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg"  example:"ok"`
	Data string `json:"data"  example:"{id: codeID, code: code, username: username, status: 1, expired: 1, expiredTime: 2025-01-01 00:00:00, deleted: 1}"`
}

type GetInviteCodeErrorResponse struct {
	Code int    `json:"code" example:"-1"`
	Msg  string `json:"msg"  example:"error msg"`
	Data string `json:"data" example:"null"`
}
