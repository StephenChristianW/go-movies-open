package UserService

import (
	"errors"
	"fmt"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/collections"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/db"
	"github.com/StephenChristianW/go-movies-open/db/RedisService/ValidateUser"
	"github.com/StephenChristianW/go-movies-open/security/SecurityBcrypt"
	"github.com/StephenChristianW/go-movies-open/utils/UtilsTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var userExpireDuration = UtilsTime.DaysToTimeDuration(config.UserExpiration)

func codeIsBind(code string) (bool, error) {
	ctx, cancel := db.GetCtx()
	defer cancel()
	// 查询是否有未删除、已使用的邀请码
	result := collections.GetInviteCodeCollection().FindOne(ctx, bson.M{
		"code":    code,
		"status":  2, // 已使用
		"deleted": 1, // 未删除
	})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, result.Err()
	}
	return true, nil
}

func newUser(user CreateUser) (*UserSchemas, error) {

	pwd, err := SecurityBcrypt.GenerateHashPwd(user.Password)
	if err != nil {
		return nil, err
	}
	return &UserSchemas{
		Username:    user.Username,
		Password:    pwd,
		Code:        user.Code,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
		CreateAtStr: UtilsTime.NowTime(),
		UpdateAtStr: UtilsTime.NowTime(),
		Status:      1,
		Deleted:     1,
		ContactType: user.ContactType,
		ContactInfo: user.ContactInfo,
	}, nil

}

// calculateLockSeconds 计算用户锁定剩余秒数，如果已解锁自动更新状态
func calculateLockSeconds(userID string, lockDuration time.Duration, lastTime time.Time) (int, error) {
	unlockTime := lastTime.Add(lockDuration)
	unlocked, remain := UtilsTime.CheckUnlockTime(unlockTime)
	objId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}
	if unlocked {
		// 已解锁，更新用户状态为正常
		ctx, cancel := db.GetCtx()
		defer cancel()
		_, err := collections.GetUserCollection().UpdateOne(ctx,
			bson.M{"_id": objId, "deleted": 1},
			bson.M{"$set": bson.M{"status": 1}})
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	return remain, nil
}

// validateUser 根据失败次数和最后失败时间判断是否锁定
func validateUser(userID string) (int, int, error) {
	has, data, err := ValidateUser.HasFailedLoginRecord(userID)
	if err != nil {
		return -1, 0, err
	}
	if !has {
		return 0, 0, nil // 没有失败记录，直接放行
	}
	var lockDuration time.Duration
	switch data.Count {
	case 3:
		lockDuration = time.Minute * 1
	case 5:
		lockDuration = time.Minute * 5
	case 7:
		lockDuration = time.Minute * 30
	case 10:
		lockDuration = time.Hour * 24
	case 11:
		// 超过阈值，永久冻结
		objId, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return 0, data.Count, err
		}
		ctx, cancel := db.GetCtx()
		defer cancel()
		_, _ = collections.GetUserCollection().UpdateOne(ctx,
			bson.M{"_id": objId, "deleted": 1},
			bson.M{"$set": bson.M{"status": 3}})
		return -1, data.Count, fmt.Errorf("密码错误:下次重置密码前用户已冻结登录")
	default:
		return 0, data.Count, nil
	}
	// 普通锁定逻辑
	seconds, err := calculateLockSeconds(userID, lockDuration, data.Last)
	if err != nil {
		return seconds, data.Count, err
	}
	if seconds > 0 {
		return seconds, data.Count, nil
	}
	return 0, data.Count, nil
}
