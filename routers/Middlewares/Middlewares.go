package Middlewares

import (
	"errors"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/collections"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	redis "github.com/StephenChristianW/go-movies-open/db/RedisService/OperateToken"
	"github.com/StephenChristianW/go-movies-open/utils/Jwt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

// ==================================== 中间件 ====================================
const (
	AdminInfo = "adminInfo"
	UserInfo  = "userInfo"
)

// AdminAuthMiddleware 管理员鉴权中间件
// - 排除 /admin/login 接口
// - 校验请求头中的 Authorization token
// - 解析 JWT 并校验 Redis 中 token 是否有效
// - 将解析出的管理员信息存入 context
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 登录接口无需鉴权 (支持版本号或前缀)
		if strings.HasSuffix(c.FullPath(), "/login") {
			c.Next()
			return
		}
		if strings.HasSuffix(c.FullPath(), "/api") {
			c.Next()
			return
		}
		if strings.HasSuffix(c.FullPath(), "/rlogin") {
			c.Next()
			return
		}

		// 从 Header 获取 token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"code": 401, "msg": "token缺失"})
			c.Abort()
			return
		}

		// 解析 JWT
		adminClaims, err := Jwt.ParseAdminToken(token)
		if err != nil {
			c.JSON(401, gin.H{"code": 401, "msg": "token无效"})
			c.Abort()
			return
		}

		// 校验 Redis 中是否存在该 token
		redisToken, err := redis.GetAdminToken(adminClaims.AdminID)
		if err != nil || redisToken != token {
			c.JSON(401, gin.H{"code": 401, "msg": "token已失效"})
			c.Abort()
			return
		}

		// 将管理员信息存入 Gin 上下文
		c.Set(AdminInfo, adminClaims)
		c.Next()
	}
}

func isIPBlocked(ip string) (bool, error) {
	ctx, cancel := db.GetCtx()
	defer cancel()
	result := collections.GetBlockIps().FindOne(ctx, bson.M{"ip": ip})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return false, nil
		}
		return true, result.Err()
	}
	return true, nil
}

// UserAuthMiddleware 用户鉴权中间件
// 用于检查请求中携带的 JWT token 是否有效，并校验 Redis 中的 token
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		blocked, err := isIPBlocked(ip) // 查询数据库
		if err != nil {
			// 查询异常，可选择放行或者阻止
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  "服务器异常",
			})
			return
		}
		if blocked {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": -5,
				"msg":  "该 IP 已被封禁",
			})
			return
		}
		c.Next() // 继续处理
		// 以下接口无需鉴权
		// 登录接口无需鉴权 (支持版本号或前缀)
		if strings.HasSuffix(c.FullPath(), "/login") ||
			strings.HasSuffix(c.FullPath(), "/register") ||
			strings.HasSuffix(c.FullPath(), "/forget") {
			c.Next() // 放行
			return
		}

		// 从请求 Header 获取 token
		token := c.GetHeader("Authorization")
		if token == "" {
			// token 缺失 → 拦截请求
			c.JSON(401, gin.H{"code": 401, "msg": "token缺失"})
			c.Abort()
			return
		}

		// 解析 JWT，提取用户信息
		userClaims, err := Jwt.ParseUserToken(token)
		if err != nil {
			// token 格式不正确或已过期 → 拦截请求
			c.JSON(401, gin.H{"code": 401, "msg": "token无效"})
			c.Abort()
			return
		}

		// 校验 Redis 中存储的 token 是否与请求 token 匹配
		// 防止 token 被非法使用或已被手动清理
		redisToken, err := redis.GetUserToken(userClaims.UserID)
		if err != nil || redisToken != token {
			c.JSON(401, gin.H{"code": 401, "msg": "token已失效"})
			c.Abort()
			return
		}

		// 将解析出的用户信息存入 Gin 上下文，供后续 handler 使用
		c.Set(UserInfo, userClaims)
		c.Next()
	}
}
