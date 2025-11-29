package OperateToken

import (
	"github.com/StephenChristianW/go-movies-open/db/RedisService"
	"time"
)

var redisClient = RedisService.RedisClient
var ctx = RedisService.Ctx

func StoreAdminToken(adminID, token string, expire time.Duration) error {
	key := "admin_token:" + adminID
	return redisClient.Set(ctx, key, token, expire).Err()
}
func GetAdminToken(adminID string) (string, error) {
	key := "admin_token:" + adminID
	token, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			return "", err
		}
		return "", nil
	}
	return token, nil
}
func DeleteAdminToken(adminID string) error {
	key := "admin_token:" + adminID
	return redisClient.Del(ctx, key).Err()
}

func StoreUserToken(userID, token string, expire time.Duration) error {
	key := "user_token:" + userID
	return redisClient.Set(ctx, key, token, expire).Err()
}
func GetUserToken(userID string) (string, error) {
	key := "user_token:" + userID
	token, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			return "", err
		}
		return "", nil
	}
	return token, nil
}
func DeleteUserToken(userID string) error {
	key := "user_token:" + userID
	return redisClient.Del(ctx, key).Err()
}
