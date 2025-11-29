package AdminService

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//======================== 管理员数据库结构 ========================

// AdminInDB 数据库存储的管理员信息结构体
// 用于直接映射 MongoDB 中的 admin 集合
type AdminInDB struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`            // 数据库唯一ID
	Username string             `json:"username" bson:"username"` // 管理员用户名
	Password string             `json:"password" bson:"password"` // 密码哈希
	Status   int                `json:"status" bson:"status"`     // 是否启用: 1-已启用 2-未启用
	Deleted  int                `json:"deleted" bson:"deleted"`   // 是否删除: 1-未删除 2-已删除
	Role     string             `json:"role" bson:"role"`         // 角色ID
}

// AdminCreate 创建管理员请求数据结构
type AdminCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Deleted  int    `json:"deleted"`
	Role     string `json:"role" bson:"role"`
}

// Admin 对外展示的管理员信息结构体
type Admin struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Status   int    `json:"status" bson:"status"`
	Deleted  int    `json:"deleted" bson:"deleted"`
	Role     string `json:"role" bson:"role"`
}

// AdminList 分页查询管理员请求结构体
type AdminList struct {
	IDs      []string `json:"ids" bson:"ids"`             // 可选，指定ID列表筛选
	Username string   `json:"username" bson:"username"`   // 可选，用户名模糊筛选
	Status   int      `json:"status" bson:"status"`       // 可选，状态筛选
	Deleted  int      `json:"deleted" bson:"deleted"`     // 可选，删除状态筛选
	Role     string   `json:"role" bson:"role"`           // 可选，角色筛选
	Page     int      `json:"page" bson:"page"`           // 当前页码
	PageSize int      `json:"page_size" bson:"page_size"` // 每页数量
}

//======================== 认证与会话 ========================

// AuthTokens 认证令牌对
// 包含访问令牌和刷新令牌，用于API访问和令牌刷新
type AuthTokens struct {
	Token  string `json:"token" bson:"token"` // 短期访问令牌
	RToken string `json:"rt" bson:"rt"`       // 长期刷新令牌
}

// DeviceSession 设备会话信息
// 记录用户在特定设备上的活跃会话状态
type DeviceSession struct {
	DeviceId   string    `json:"device_id" bson:"device_id"`     // 设备会话ID
	RToken     string    `json:"rt" bson:"rt"`                   // 当前会话刷新令牌
	AuthStatus int       `json:"auth_status" bson:"auth_status"` // 记住登录状态: 1-记住 2-不记住
	IP         string    `json:"ip" bson:"ip"`                   // 设备当前IP
	LastActive time.Time `json:"last_active" bson:"last_active"` // 最后活跃时间
}

// UserAuthProfile 用户认证档案
// 记录用户所有设备会话和认证配置信息
type UserAuthProfile struct {
	UserID      string                   `json:"user_id" bson:"user_id"`           // 用户ID
	Sessions    map[string]DeviceSession `json:"sessions" bson:"sessions"`         // 活跃会话集合，key 为设备ID
	MaxSessions int                      `json:"max_sessions" bson:"max_sessions"` // 最大并发会话数
}
