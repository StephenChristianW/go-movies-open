package AdminService

import (
	"errors"
	"fmt"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/collections"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	"github.com/StephenChristianW/go-movies-open/db/RedisService/OperateToken"
	"github.com/StephenChristianW/go-movies-open/security/SecurityBcrypt"
	"github.com/StephenChristianW/go-movies-open/utils/Jwt"
	"github.com/StephenChristianW/go-movies-open/utils/UtilsTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// =======================================
//
//	Token 有效期配置
//
// =======================================
var adminTokenDuration = UtilsTime.DaysToTimeDuration(config.AdminExpiration)
var adminRefreshTokenDuration = UtilsTime.DaysToTimeDuration(config.AdminRExpiration)

// =======================================
//          登录与记住密码
// =======================================

// LoginAndRemember 支持两种登录方式：
// 1. 使用用户名+密码登录
// 2. 使用 rToken 免密登录
func (*AdminService) LoginAndRemember(loginData LoginCredentialsRequest) (AuthTokens, error) {
	var tokens AuthTokens
	// 密码登录
	if loginData.RToken == "" || loginData.Password != "" {
		admin, err := verifyAdminByPassword(loginData.Username, loginData.Password)
		if err != nil {
			return tokens, err
		}
		return generateAuthTokens(admin.ID, loginData.DeviceId, admin.Username, loginData.IP, loginData.AuthStatus)
	}

	// rToken 免密登录
	if loginData.RToken != "" && loginData.Password == "" {
		admin, err := verifyAdminByRToken(loginData)
		if err != nil {
			return tokens, err
		}
		return generateAuthTokens(admin.ID, loginData.DeviceId, admin.Username, loginData.IP, loginData.AuthStatus)
	}

	return tokens, errors.New("登录信息错误")
}

// AdminLogin 使用用户名+密码登录并返回 JWT
func (*AdminService) AdminLogin(username, password, ip, deviceID string) (string, error) {

	admin, err := verifyAdminByPassword(username, password)
	if err != nil {
		return "", err
	}

	// TODO 完善后台控制AuthStatus

	deleteDeviceSessionSafe(admin.ID, deviceID, ip)

	return getOrCreateAdminToken(admin.ID, username)
}

// LogOut 管理员登出，删除 Redis 中 token
func (*AdminService) LogOut(adminID string) error {
	log.Printf("管理员 %s 已登出\n", adminID)
	return deleteAdminToken(adminID)
}

// =======================================
//          Token / rToken 相关
// =======================================

// generateAuthTokens 生成 JWT 与 rToken，并保存 rToken
func generateAuthTokens(adminID, deviceID, username, ip string, authStatus int) (AuthTokens, error) {
	var tokens AuthTokens

	// 生成 rToken
	rt, err := Jwt.GeneratedAdminRToken(adminID, deviceID, ip, time.Now().UTC(), adminRefreshTokenDuration)
	if err != nil {
		return tokens, err
	}

	// 生成 JWT
	token, err := getOrCreateAdminToken(adminID, username)
	if err != nil {
		return tokens, err
	}

	// 更新 rToken 到 MongoDB
	if err := saveRefreshToken(adminID, rt, ip, deviceID, authStatus); err != nil {
		return tokens, err
	}

	tokens.Token = token
	tokens.RToken = rt
	return tokens, nil
}

// saveRefreshToken 更新或插入 rToken 到 MongoDB（upsert）
func saveRefreshToken(adminID, rToken, ip, deviceID string, authStatus int) error {
	ctx, cancel := db.GetCtx()
	defer cancel()

	_, err := collections.GetAdminRememberCollection().UpdateOne(
		ctx,
		bson.M{"user_id": adminID},
		bson.M{
			"$set": bson.M{
				"sessions." + deviceID: DeviceSession{
					DeviceId:   deviceID,
					AuthStatus: authStatus,
					RToken:     rToken,
					IP:         ip,
					LastActive: time.Now().UTC(),
				},
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}
func deleteDeviceSessionSafe(userID, deviceID, ip string) {
	ctx, cancel := db.GetCtx()
	defer cancel()

	// 直接尝试删除，不存在也没关系
	_, _ = collections.GetAdminRememberCollection().UpdateOne(
		ctx,
		bson.M{
			"user_id":                      userID,
			"sessions." + deviceID + ".ip": ip, // IP 匹配
		},
		bson.M{
			"$unset": bson.M{
				"sessions." + deviceID: "",
			},
		},
	)
	// 不管有没有匹配到都不报错
}

// getOrCreateAdminToken 获取 Redis token，如果不存在则生成新的 JWT
func getOrCreateAdminToken(adminID, username string) (string, error) {
	token, err := OperateToken.GetAdminToken(adminID)
	if err != nil {
		return "", err
	}

	if token == "" {
		token, err = Jwt.GenerateAdminToken(adminID, username, adminTokenDuration)
		if err != nil {
			return "", err
		}
		if err = OperateToken.StoreAdminToken(adminID, token, adminTokenDuration); err != nil {
			return "", err
		}
	}

	return token, nil
}

// deleteAdminToken 删除 Redis 中的管理员 token
func deleteAdminToken(adminID string) error {
	token, err := OperateToken.GetAdminToken(adminID)
	if err != nil {
		return fmt.Errorf("获取管理员 Redis token 失败: %v", err)
	}
	if token != "" {
		if err := OperateToken.DeleteAdminToken(adminID); err != nil {
			return fmt.Errorf("删除管理员 Redis token 失败: %v", err)
		}
	}
	return nil
}

// =======================================
//          管理员验证
// =======================================

// verifyAdminByPassword 校验用户名和密码
func verifyAdminByPassword(username, password string) (Admin, error) {
	ctx, cancel := db.GetCtx()
	defer cancel()

	var admin AdminInDB
	err := collections.GetAdminCollection().FindOne(ctx, bson.M{
		"username": username,
		"deleted":  1,
	}).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Admin{}, fmt.Errorf("用户名错误或管理员 %s 已被禁用", username)
		}
		return Admin{}, err
	}

	if admin.Status == 2 {
		return Admin{}, errors.New("该账号已被禁用")
	}

	if !SecurityBcrypt.CompareHashPwd(admin.Password, password) {
		return Admin{}, errors.New("密码错误")
	}

	return Admin{
		ID:       admin.ID.Hex(),
		Username: admin.Username,
		Role:     admin.Role,
		Status:   admin.Status,
		Deleted:  admin.Deleted,
	}, nil
}

// verifyAdminByRToken 校验 rToken，返回管理员信息
func verifyAdminByRToken(loginData LoginCredentialsRequest) (Admin, error) {
	ctx, cancel := db.GetCtx()
	defer cancel()

	var result Admin

	rtClaims, err := Jwt.ParseAdminRToken(loginData.RToken)
	if err != nil {
		return result, err
	}
	if loginData.DeviceId != rtClaims.DeviceID {
		return result, errors.New("设备ID不一致")
	}

	if loginData.IP != rtClaims.IP {
		return result, errors.New("设备IP不一致")
	}

	var userAuth UserAuthProfile
	err = collections.GetAdminRememberCollection().FindOne(
		ctx,
		bson.M{
			"user_id": rtClaims.AdminID,
			"sessions." + rtClaims.DeviceID + ".device_id": rtClaims.DeviceID,
			"sessions." + rtClaims.DeviceID + ".ip":        rtClaims.IP,
		},
	).Decode(&userAuth)
	if err != nil {
		return result, err
	}
	var admin AdminInDB
	err = collections.GetAdminCollection().FindOne(ctx, bson.M{"username": loginData.Username}).Decode(&admin)
	if err != nil {
		return result, err
	}

	if admin.ID.Hex() != userAuth.UserID {
		return result, errors.New("用户不匹配")
	}
	result.ID = admin.ID.Hex()
	result.Username = admin.Username
	result.Role = admin.Role
	result.Status = admin.Status
	result.Deleted = admin.Deleted

	return result, nil
}
