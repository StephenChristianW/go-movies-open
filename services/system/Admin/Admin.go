package AdminService

import "github.com/StephenChristianW/go-movies-open/services"

// ======================= AdminService 接口定义 =======================

// AdminService 实现管理员业务逻辑的服务结构体
type AdminService struct{}

// AdminInterface 定义管理员服务应提供的方法
type AdminInterface interface {

	// AdminLogin 管理员登录，返回短期访问令牌
	AdminLogin(username, password, ip, deviceID string) (token string, err error)

	// LoginAndRemember 登录并生成访问令牌 + 刷新令牌
	LoginAndRemember(data LoginCredentialsRequest) (AuthTokens, error)

	// LogOut 管理员登出
	LogOut(adminId string) error

	// AdminList 获取管理员列表，支持分页和条件过滤
	AdminList(filter AdminList) (services.Pagination, error)
	// CreateAdmin 创建新管理员
	CreateAdmin(username, password, roleID string) error

	// DeleteAdmin 逻辑删除管理员
	DeleteAdmin(nowUsername, delUsername string) error

	// ChangeAdminPassword 修改管理员密码
	ChangeAdminPassword(username, newPassword string) error

	// BannedAdminUser 禁用管理员
	BannedAdminUser(bannedUsername string) error

	// ActiveAdminUser 启用管理员
	ActiveAdminUser(bannedUsername string) error
}
