package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecret 生成指定长度的随机 secret
func GenerateSecret(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// 用 base64 编码，方便存储和使用
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func main() {
	fmt.Println(GenerateSecret(32))

}
