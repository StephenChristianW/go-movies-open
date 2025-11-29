package UserService

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserFilter struct {
	IDs         []string `json:"ids" bson:"_ids"`
	UserName    string   `bson:"username" json:"username"`
	Status      int      `bson:"status" json:"status"`
	Deleted     int      `json:"deleted" bson:"deleted"`           // 是否删除 1 未删除 2 已删除
	ContactType string   `json:"contact_type" bson:"contact_type"` // 联系渠道-用于找回密码
	ContactInfo string   `json:"contact_info" bson:"contact_info"` // 联系信息-用于找回密码
	BeginTime   string   `json:"begin_time" bson:"begin_time"`
	EndTime     string   `json:"end_time" bson:"end_time"`
}
type UserSchemas struct {
	ID          string    `json:"id" bson:"_id"`
	Username    string    `json:"username" bson:"username"`
	Password    string    `json:"password" bson:"password"`
	Code        string    `json:"code" bson:"code"`           // 邀请码
	CreateAt    time.Time `json:"create_at" bson:"create_at"` // 创建时间
	CreateAtStr string    `json:"create_at_str" bson:"create_at_str"`
	UpdateAt    time.Time `json:"update_at" bson:"update_at"` // 更新时间
	UpdateAtStr string    `json:"update_at_str" bson:"update_at_str"`
	Status      int       `json:"status" bson:"status"`             // 是否禁用 1 正常 2 已禁用
	Deleted     int       `json:"deleted" bson:"deleted"`           // 是否删除 1 未删除 2 已删除
	ContactType string    `json:"contact_type" bson:"contact_type"` // 联系渠道-用于找回密码
	ContactInfo string    `json:"contact_info" bson:"contact_info"` // 联系信息-用于找回密码
}
type UserInDB struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Username    string             `json:"username" bson:"username"`
	Password    string             `json:"password" bson:"password"`
	Code        string             `json:"code" bson:"code"`           // 邀请码
	CreateAt    time.Time          `json:"create_at" bson:"create_at"` // 创建时间
	CreateAtStr string             `json:"create_at_str" bson:"create_at_str"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"` // 更新时间
	UpdateAtStr string             `json:"update_at_str" bson:"update_at_str"`
	Status      int                `json:"status" bson:"status"`             // 是否禁用 1 正常 2 已禁用
	Deleted     int                `json:"deleted" bson:"deleted"`           // 是否删除 1 未删除 2 已删除
	ContactType string             `json:"contact_type" bson:"contact_type"` // 联系渠道-用于找回密码
	ContactInfo string             `json:"contact_info" bson:"contact_info"` // 联系信息-用于找回密码
}
type CreateUser struct {
	Username    string `json:"username" binding:"required" bson:"username"`
	Password    string `json:"password"  binding:"required" bson:"password"`
	Code        string `json:"code"  binding:"required" bson:"code"`
	ContactType string `json:"contact_type" bson:"contact_type"` // 联系渠道-用于找回密码
	ContactInfo string `json:"contact_info" bson:"contact_info"` // 联系信息-用于找回密码
}
type NewUserSchemas struct {
	Username    string    `json:"username" binding:"required" bson:"username"`
	Password    string    `json:"password" binding:"required" bson:"password"`
	Code        string    `json:"code" binding:"required" bson:"code"` // 邀请码
	CreateAt    time.Time `json:"create_at" bson:"create_at"`          // 创建时间
	CreateAtStr string    `json:"create_at_str" bson:"create_at_str"`
	UpdateAt    time.Time `json:"update_at" bson:"update_at"` // 更新时间
	UpdateAtStr string    `json:"update_at_str" bson:"update_at_str"`
	Status      int       `json:"status" bson:"status"`             // 是否禁用 1 正常 2 已禁用禁止一切操作 3 已冻结 可更改密码解锁登录
	Deleted     int       `json:"deleted" bson:"deleted"`           // 是否删除 1 未删除 2 已删除
	ContactType string    `json:"contact_type" bson:"contact_type"` // 联系渠道-用于找回密码
	ContactInfo string    `json:"contact_info" bson:"contact_info"` // 联系信息-用于找回密码
}
type UserLogin struct {
	Username string `json:"username" binding:"required" bson:"username"`
	Password string `json:"password" binding:"required" bson:"password"`
}

type ForgetUser struct {
	Username string `json:"username" binding:"required" bson:"username"`
	Code     string `json:"code" binding:"required" bson:"code"`
	Password string `json:"password" binding:"required" bson:"password"`
}

type ChangePasswordSchema struct {
	OldPassword string `json:"old_password" binding:"required" bson:"old_password"`
	NewPassword string `json:"new_password" binding:"required" bson:"new_password"`
}

type LoginFailureResponse struct {
	FailedAttempts int `json:"num"`
	LockSeconds    int `json:"seconds"`
}

type LoginSuccessResponse struct {
	Token string `json:"token"`
}

type LoginResponse struct {
	FailedAttempts int    `json:"num"`
	LockSeconds    int    `json:"seconds"`
	Token          string `json:"token"`
}

type Rename struct {
	Username string `json:"username" bson:"username"`
}
