package config

// ======================= 全局环境常量 =======================
// LOCAL: 本地开发环境
// PROD: 生产环境
const (
	LOCAL  = "local"
	PROD   = "production"
	envLog = "[ENV] " // 日志打印前缀
)

// ======================= 数据库集合常量 =======================
// 用于统一管理 MongoDB 集合名称
const (
	InviteCodeColl = "invite_code"    // 邀请码集合，code字段为唯一索引
	AdminColl      = "admin"          // 管理员集合
	UserColl       = "user"           // 普通用户集合
	BlockIps       = "block_ips"      // 被封禁 IP 集合
	AdminRemember  = "admin_remember" // 管理员设备会话集合
)

var (
	// ======================= 服务端口 =======================

	Addr = getAddr() // 服务监听端口，默认从 ENV 获取
	Port = getPort()
	// ======================= JWT 配置 =======================

	AdminSecret      = getenv("ADMIN_SECRET", "")   // 管理员访问令牌密钥
	UserSecret       = getenv("USER_SECRET", "")    // 普通用户访问令牌密钥
	AdminRSecret     = getenv("ADMIN_R_SECRET", "") // 管理员刷新令牌密钥
	AdminExpiration  = getAdminExpirationDay()      // 管理员 token 过期天数
	AdminRExpiration = getAdminRExpirationDay()     // 管理员记住密码 token 过期天数
	UserExpiration   = getUserExpirationDay()       // 普通用户 token 过期天数

	// ======================= 数据库配置 =======================

	DBName = getDBName() // MongoDB 数据库名称
	DBUrl  = getDBUrl()  // MongoDB 连接 URL

	// ======================= Redis 配置 =======================

	RedisAddr     = getRedisAddr()     // Redis 地址，格式 ip:port
	RedisPassword = getRedisPassword() // Redis 密码
	RedisDB       = getRedisDB()       // Redis 数据库编号

	// ======================= 分页 & 批量操作配置 =======================

	MaxCreateBatchSize = getEnvInt("MAX_CREATE_BATCH_SIZE", 50000) // 最大一次批量创建数量
	DefaultPageSize    = getEnvInt("DEFAULT_PAGE_SIZE", 20)        // 默认分页条数
	MinPageSize        = getEnvInt("MIN_PAGE_SIZE", 10)            // 最小分页条数
	MaxPageSize        = getEnvInt("MAX_PAGE_SIZE", 100)           // 最大分页条数

	// ======================= 初始管理员 & 字段长度限制 =======================

	RootAdmins       = getRootAdmins()                // 初始管理员列表，从 ENV 获取
	LessLenAdminName = getEnvInt("LESS_LEN_ADMIN", 3) // 管理员用户名最小长度
	LessLenUserName  = getEnvInt("LESS_LEN_USER", 3)  // 普通用户用户名最小长度
	LessLenPwd       = getEnvInt("LESS_LEN_PWD", 6)   // 密码最小长度
)
