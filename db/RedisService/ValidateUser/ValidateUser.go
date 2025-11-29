package ValidateUser

import (
	"errors"
	"github.com/StephenChristianW/go-movies-open/db/RedisService"
	"strconv"
	"time"
)

var redisClient = RedisService.RedisClient
var ctx = RedisService.Ctx

type FailedLoginInfo struct {
	Count int       `json:"count"`
	Last  time.Time `json:"last"`
}

func SetFailedLoginInfo(userID string, info FailedLoginInfo) error {
	key := "login:failed:" + userID
	data := map[string]interface{}{
		"count": info.Count,
		"last":  info.Last.Unix(),
	}
	return redisClient.HSet(ctx, key, data).Err()
}

func getFailedLoginInfo(userID string) (FailedLoginInfo, error) {
	key := "login:failed:" + userID
	var info FailedLoginInfo
	vals, err := redisClient.HMGet(ctx, key, "count", "last").Result()
	if err != nil {
		return info, err
	}
	if len(vals) != 2 {
		return info, errors.New("数据错误")
	}
	// count
	if vals[0] != nil {
		switch v := vals[0].(type) {
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				info.Count = i
			}
		}
	}
	// last
	if vals[1] != nil {
		switch v := vals[1].(type) {
		case string:
			if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
				info.Last = time.Unix(ts, 0)
			}
		}
	}
	return info, nil
}

func HasFailedLoginRecord(userID string) (bool, FailedLoginInfo, error) {
	info, err := getFailedLoginInfo(userID)
	if err != nil {
		// Redis 异常或数据错误，可以当作没有记录
		return false, info, err
	}

	// 判断是否有有效失败记录
	if info.Count == 0 && info.Last.IsZero() {
		return false, info, nil
	}

	return true, info, nil
}
func DeleteFailedLoginRecord(userID string) error {
	key := "login:failed:" + userID
	return redisClient.Del(ctx, key).Err()
}
