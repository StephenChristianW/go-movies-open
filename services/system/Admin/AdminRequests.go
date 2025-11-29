package AdminService

//======================== controller 请求 DTO ========================

// AdminLoginRequest 登录请求结构体
// 用于前端登录请求，包含用户名、密码、设备信息和是否记住登录
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`  // 登录用户名，必填
	Password string `json:"password" binding:"required"`  // 登录密码，必填
	DeviceId string `json:"device_id" binding:"required"` // 设备唯一标识，用于多设备管理
	IP       string `json:"ip" binding:"required"`        // 当前登录IP地址
}

// BannedUserRequest 禁用管理员请求结构体
// 仅需传入待禁用的用户名
type BannedUserRequest struct {
	BannedUsername string `json:"username" binding:"required"` // 待禁用的管理员用户名
}

// DeletedUserRequest 删除管理员请求结构体
// 仅需传入待删除的用户名
type DeletedUserRequest struct {
	DelUsername string `json:"username" binding:"required"` // 待删除的管理员用户名
}

// ChangeAdminPasswordRequest 修改管理员密码请求结构体
type ChangeAdminPasswordRequest struct {
	Username    string `json:"username" binding:"required"`     // 管理员用户名
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

// LoginCredentialsRequest 用户登录凭证
// 用于后端认证逻辑，包含所有必需的登录信息和设备信息
type LoginCredentialsRequest struct {
	Username   string `json:"username" bson:"username"`       // 登录用户名，唯一标识
	Password   string `json:"password" bson:"password"`       // 登录明文密码，由后端加密验证
	DeviceId   string `json:"device_id" bson:"device_id"`     // 设备唯一标识，用于多设备会话管理
	RToken     string `json:"rt" bson:"rt"`                   // 刷新令牌(Refresh Token)，用于获取新访问令牌
	AuthStatus int    `json:"auth_status" bson:"auth_status"` // 记住登录状态：1-记住 2-不记住
	IP         string `json:"ip" bson:"ip"`                   // 登录IP地址，用于安全审计
}
