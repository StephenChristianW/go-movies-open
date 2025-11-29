package InviteCodeService

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateInviteCode struct {
	Code        string `json:"code" bson:"code"`                // 邀请码
	Username    string `json:"username" bson:"username"`        // 绑定的用户名，""为未绑定
	Status      int    `json:"status" bson:"status"`            // 使用状态 1 已使用 2 未过期
	IsExpired   int    `json:"expired" bson:"is_expired"`       // 是否过期 1 未过期 2 已过期
	ExpiredTime string `json:"expiredTime" bson:"expired_time"` // 过期时间字符串
	Deleted     int    `json:"deleted" bson:"deleted"`          // 是否删除 1 未删除 2 已删除
}

type UpdateInviteCode struct {
	Ids    []string `json:"ids" bson:"_ids"`
	Status int      `json:"status" bson:"status"` // 使用状态 1 已使用 2 未过期
	Days   int      `json:"days" bson:"days"`     // 过期延后天数
}

type CreateInviteFilter struct {
	IDs       []string `json:"ids" bson:"_ids"`
	Code      string   `json:"code" bson:"code"`          // 邀请码
	Username  string   `json:"username" bson:"username"`  // 绑定的用户名，""为未绑定
	Status    int      `json:"status" bson:"status"`      // 使用状态 1 已使用 2 未过期
	IsExpired int      `json:"expired" bson:"is_expired"` // 是否过期 1 未过期 2 已过期
	Deleted   int      `json:"deleted" bson:"deleted"`    // 是否删除 1 未删除 2 已删除

	BeginTime string `json:"beginTime" bson:"beginTime"`
	EndTime   string `json:"endTime" bson:"endTime"`
	Page      int    `json:"page" bson:"page"`
	PageSize  int    `json:"pageSize" bson:"pageSize"`
}

type InviteCodeInDB struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Code        string             `json:"code" bson:"code" `              // 邀请码
	Username    string             `json:"username" bson:"username"`       // 绑定的用户名，""为未绑定
	Status      int                `json:"status" bson:"status"`           // 使用状态 1 已使用 2 未过期
	IsExpired   int                `json:"expired" bson:"isExpired"`       // 是否过期 1 未过期 2 已过期
	ExpiredTime string             `json:"expiredTime" bson:"expiredTime"` // 过期时间字符串
	Deleted     int                `json:"deleted" bson:"deleted"`         // 是否删除 1 未删除 2 已删除
}
type InviteCode struct {
	ID          string `json:"id" bson:"_id"`
	Code        string `json:"code" bson:"code" `              // 邀请码
	Username    string `json:"username" bson:"username"`       // 绑定的用户名，""为未绑定
	Status      int    `json:"status" bson:"status"`           // 使用状态 1 已使用 2 未过期
	IsExpired   int    `json:"expired" bson:"isExpired"`       // 是否过期 1 未过期 2 已过期
	ExpiredTime string `json:"expiredTime" bson:"expiredTime"` // 过期时间字符串
	Deleted     int    `json:"deleted" bson:"deleted"`         // 是否删除 1 未删除 2 已删除
}

// InviteCodeRequest 获取邀请码请求体
type InviteCodeRequest struct {
	Code string `form:"code" binding:"required"`
}
