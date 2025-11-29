package Jwt

import (
	"fmt"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AdminClaims struct {
	AdminID              string `json:"id" bson:"_id"`
	AdminUsername        string `json:"username" bson:"username"`
	jwt.RegisteredClaims `json:"claims" bson:"claims"`
}
type AdminRClaims struct {
	AdminID    string    `json:"id" bson:"id"`
	DeviceID   string    `json:"device_id" bson:"device_id"`
	IP         string    `json:"ip" bson:"ip"`
	LastActive time.Time `json:"last_login" bson:"last_login"`
	jwt.RegisteredClaims
}

func GeneratedAdminRToken(adminId, deviceID, ip string, lastActive time.Time, expireDuration time.Duration) (string, error) {
	claims := AdminRClaims{
		AdminID:    adminId,
		DeviceID:   deviceID,
		IP:         ip,
		LastActive: lastActive,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AdminRSecret))
}

func ParseAdminRToken(tokenStr string) (*AdminRClaims, error) {
	claims := &AdminRClaims{}
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AdminRSecret), nil
	})
	if err != nil {
		return nil, err
	}
	// 检查 token 是否有效
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func GenerateAdminToken(adminId, adminUsername string, expireDuration time.Duration) (string, error) {
	claims := AdminClaims{
		AdminID:       adminId,
		AdminUsername: adminUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AdminSecret))
}
func ParseAdminToken(tokenStr string) (*AdminClaims, error) {
	claims := &AdminClaims{}
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AdminSecret), nil
	})
	if err != nil {
		return nil, err
	}
	// 检查 token 是否有效
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

type UserClaims struct {
	UserID   string `json:"user_id" bson:"user_id"`
	Username string `json:"username"  bson:"username"`
	jwt.RegisteredClaims
}

func GenerateUserToken(userID, userName string, expireDuration time.Duration) (string, error) {
	claims := UserClaims{
		UserID:   userID,
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.UserSecret))
}

func ParseUserToken(tokenStr string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.UserSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
