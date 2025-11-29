package UserService

import (
	"errors"
	"fmt"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/collections"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	"github.com/StephenChristianW/go-movies-open/db/RedisService/OperateToken"
	"github.com/StephenChristianW/go-movies-open/db/RedisService/ValidateUser"
	"github.com/StephenChristianW/go-movies-open/security/SecurityBcrypt"
	"github.com/StephenChristianW/go-movies-open/services"
	"github.com/StephenChristianW/go-movies-open/services/System/InviteCode"
	"github.com/StephenChristianW/go-movies-open/utils/Jwt"
	"github.com/StephenChristianW/go-movies-open/utils/UtilsTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserInterface interface {
	RegisterUser(user CreateUser) error
	UserLogin(username, password string) (LoginResponse, error)
	ForgetPassword(username, code, password string) error
	ChangePassword(userID, oldPassword, newPassword string) error
	RenameUser(userID, newUserName string) error
	DeleteUser(username string) error
	BannedUser(username string) error
	UserList(filter UserFilter) (services.Pagination, error)
}
type UserSchema struct {
}

// RegisterUser 注册新用户
// 功能：
// 1. 检查用户名是否已存在
// 2. 检查邀请码是否已被使用
// 3. 创建新的用户数据
// 4. 将新用户插入 MongoDB
// 5. 更新邀请码状态为已绑定
func (*UserSchema) RegisterUser(user CreateUser) error {
	ctx, cancel := db.GetCtx() // 获取 MongoDB 上下文
	defer cancel()
	if len(user.Username) < config.LessLenUserName {
		return errors.New("用户名至少3个字符")
	}
	if len(user.Password) < config.LessLenPwd {
		return errors.New("密码至少6位")
	}
	// 检查用户是否存在
	var userInDB UserInDB
	err := collections.GetUserCollection().FindOne(ctx, bson.M{
		"username": user.Username,
		"deleted":  1, // 只检查未删除的用户
	}).Decode(&userInDB)
	if err == nil {
		return errors.New("用户已存在") // 用户名已存在
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return err // 查询出错
	}

	// 检查邀请码是否已被使用
	ok, err := codeIsBind(user.Code)
	if err != nil {
		return err // 查询邀请码出错
	}
	if ok {
		return fmt.Errorf("邀请码 %s 已被使用", user.Code)
	}

	// 创建新用户数据（newUser 会生成初始状态、密码哈希等）
	newU, err := newUser(user)
	if err != nil {
		return err
	}

	// 插入新用户到 MongoDB
	_, err = collections.GetUserCollection().InsertOne(ctx, newU)
	if err != nil {
		return err
	}

	// 更新邀请码状态为已绑定，并记录绑定的用户名
	_, err = collections.GetInviteCodeCollection().UpdateOne(
		ctx,
		bson.M{"code": user.Code},
		bson.M{"$set": bson.M{"username": user.Username, "status": 2}},
	)
	if err != nil {
		return err
	}

	return nil
}

// UserLogin 用户登录
// 功能：
// 1. 查询用户信息
// 2. 检查用户状态（正常、禁用、冻结）
// 3. 根据失败次数判断是否锁定
// 4. 验证密码，如果失败 记录失败次数并更新 Redis
// 5. 登录成功，清理失败记录
// 6. 获取已有 token 或生成新的 JWT token 并存储
func (*UserSchema) UserLogin(username, password string) (LoginResponse, error) {
	ctx, cancel := db.GetCtx() // 获取 MongoDB 上下文
	defer cancel()
	var resp = LoginResponse{}
	// 1. 查询用户
	var userInDB UserInDB
	err := collections.GetUserCollection().FindOne(ctx, bson.M{"username": username, "deleted": 1}).Decode(&userInDB)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return resp, errors.New("用户不存在")
		}
		return resp, err
	}
	userID := userInDB.ID.Hex() // 将 ObjectID 转成字符串

	// 2. 状态检查
	switch userInDB.Status {
	case 2:
		return resp, errors.New("账户已被禁用")
	case 3:
		return resp, errors.New("账户已被安全冻结,请重置密码后重试")
	}
	if !SecurityBcrypt.CompareHashPwd(userInDB.Password, password) {
		_ = OperateToken.DeleteUserToken(userID)
		// 3. 检查是否锁定
		// validateUser 根据 Redis 失败次数和上次失败时间判断是否锁定
		lockSeconds, count, err := validateUser(userID)
		if err != nil {
			return resp, err
		}
		if count == 0 {
			count = 1
		}
		resp.LockSeconds = lockSeconds
		resp.FailedAttempts = count
		if lockSeconds > 0 {
			return resp, nil
		}
		count += 1
		err = ValidateUser.SetFailedLoginInfo(userID,
			ValidateUser.FailedLoginInfo{
				Count: count,
				Last:  time.Now(),
			})
		if err != nil {
			return resp, err
		}
		return resp, nil
	}
	err = ValidateUser.DeleteFailedLoginRecord(userID)
	if err != nil {
		return resp, err
	}
	// 生成token
	var token string
	token, err = OperateToken.GetUserToken(userID)
	if err != nil {
		return resp, err
	}
	if token != "" {
		resp.Token = token
		return resp, nil
	}
	token, err = Jwt.GenerateUserToken(userID, username, userExpireDuration)
	if err != nil {
		return resp, err
	}
	err = OperateToken.StoreUserToken(userID, token, userExpireDuration)
	if err != nil {
		return resp, err
	}
	resp.Token = token

	return resp, nil
}

func (*UserSchema) RenameUser(userID, newUserName string) error {
	ctx, cancel := db.GetCtx()
	defer cancel()
	objId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	// 检查新用户名是否已被占用（排除自己）
	count, err := collections.GetUserCollection().CountDocuments(ctx, bson.M{
		"username": newUserName,
		"deleted":  1,
		"_id":      bson.M{"$ne": objId},
	})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}
	// 更新用户
	res, err := collections.GetUserCollection().UpdateOne(ctx,
		bson.M{"_id": objId},
		bson.M{"$set": bson.M{"username": newUserName}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func (*UserSchema) DeleteUser(username string) error {
	return nil
}
func (*UserSchema) BannedUser(username string) error {
	return nil
}
func filterFunc(filter UserFilter) bson.M {
	var condition bson.M
	if len(filter.IDs) > 0 {

	}
	if filter.UserName != "" {
		condition["username"] = bson.M{"$regex": filter.UserName}

	}
	return condition
}
func (*UserSchema) UserList(filter UserFilter) (services.Pagination, error) {

	return services.Pagination{}, nil
}

// ChangePassword 修改用户密码
// 参数:
//
//	userID - 用户 ID (string)
//	oldPassword - 旧密码
//	newPassword - 新密码
//
// 返回值:
//
//	error - 错误信息，如果修改失败会返回具体错误
func (*UserSchema) ChangePassword(userID, oldPassword, newPassword string) error {
	// 1. 获取 MongoDB 上下文
	ctx, cancel := db.GetCtx()
	defer cancel()

	// 2. 将字符串 userID 转换为 MongoDB ObjectID
	objId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// 3. 查询数据库获取用户信息
	var userInDB UserInDB
	err = collections.GetUserCollection().FindOne(ctx, bson.M{
		"_id":     objId,
		"deleted": 1, // 未删除状态
	}).Decode(&userInDB)
	if err != nil {
		return err
	}

	// 4. 校验旧密码是否正确
	if !SecurityBcrypt.CompareHashPwd(userInDB.Password, oldPassword) {
		return errors.New("密码错误")
	}

	// 5. 生成新密码的哈希值
	newHash, err := SecurityBcrypt.GenerateHashPwd(newPassword)
	if err != nil {
		return err
	}

	// 6. 更新数据库中的密码和更新时间
	_, err = collections.GetUserCollection().UpdateOne(ctx, bson.M{"_id": userInDB.ID},
		bson.M{"$set": bson.M{
			"password":      newHash,             // 更新密码
			"update_at":     time.Now(),          // 更新时间
			"update_at_str": UtilsTime.NowTime(), // 更新时间字符串表示
		}})
	if err != nil {
		return err
	}

	return nil
}

// ForgetPassword 忘记密码重置
// 参数:
//
//	username - 用户名
//	code - 用户的邀请码
//	password - 新密码
//
// 返回值:
//
//	error - 错误信息，如果修改失败会返回具体错误
func (*UserSchema) ForgetPassword(username, code, password string) error {
	// 1. 获取 MongoDB 上下文
	ctx, cancel := db.GetCtx()
	defer cancel()

	// 2. 查询用户信息
	var userInDB UserInDB
	err := collections.GetUserCollection().FindOne(ctx, bson.M{"username": username, "deleted": 1}).Decode(&userInDB)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("无用户: " + username)
		}
		return err
	}
	if userInDB.Status == 2 {
		return errors.New("用户已禁用, 请联系管理员")
	}
	// 3. 查询用户对应的邀请码信息
	var inviteCode InviteCodeService.InviteCodeInDB
	err = collections.GetInviteCodeCollection().FindOne(ctx, bson.M{"code": userInDB.Code}).Decode(&inviteCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("邀请码错误")
		}
		return err
	}

	// 4. 校验用户提供的邀请码是否匹配
	if code != inviteCode.Code {
		return errors.New("邀请码错误, 密码无法修改")
	}

	// 5. 生成新密码的哈希值
	pwd, err := SecurityBcrypt.GenerateHashPwd(password)
	if err != nil {
		return err
	}

	// 6. 更新用户密码及更新时间
	_, err = collections.GetUserCollection().UpdateOne(ctx, bson.M{"_id": userInDB.ID}, bson.M{"$set": bson.M{
		"status":        1,
		"password":      pwd,                 // 新密码
		"update_at":     time.Now(),          // 更新时间
		"update_at_str": UtilsTime.NowTime(), // 更新时间字符串
	}})
	_ = ValidateUser.DeleteFailedLoginRecord(userInDB.ID.Hex())
	if err != nil {
		return err
	}
	return nil
}
