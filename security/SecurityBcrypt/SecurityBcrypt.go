package SecurityBcrypt

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// GenerateHashPwd 对明文密码生成 bcrypt 哈希
func GenerateHashPwd(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成密码哈希失败:", err)
	}
	return string(hash), err
}

// CompareHashPwd 验证明文密码是否与哈希匹配
func CompareHashPwd(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//secret, err := Jwt.GenerateSecret(32) // 32字节 ≈ 43字符
//if err != nil {
//	panic(err)
//}
//fmt.Println("Generated Secret:", secret)
