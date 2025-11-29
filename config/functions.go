package config

import (
	"fmt"
	"github.com/StephenChristianW/go-movies-open/utils"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var once sync.Once
var Env = getEnv()

func ProdEnv() bool {
	switch Env {
	case PROD:
		return true
	case LOCAL:
		return false
	default:
		panic("Unknown Environment: " + Env)
	}
}

// ======================= 环境初始化 =======================

func loadEnv() {
	once.Do(func() {
		cwd, _ := os.Getwd()
		log.Println(envLog+"当前工作目录:", cwd)

		// 尝试加载默认 .env
		if err := godotenv.Load(".env"); err == nil {
			log.Println(envLog + "环境配置文件加载成功")
		} else {
			altPath := cwd + "/.env"
			if err := godotenv.Load(altPath); err == nil {
				fmt.Println("从备用路径加载 .env 成功:", altPath)
			} else {
				fmt.Println("未找到 .env 文件，使用系统环境变量")
			}
		}

		// 如果 ENV 为空，默认 local
		if os.Getenv("ENV") == "" {
			fmt.Println("ENV 未设置，默认使用 local")
			_ = os.Setenv("ENV", LOCAL)
		}
		log.Println(envLog+"当前 ENV:", os.Getenv("ENV"))
	})
}

// ======================= 获取环境变量 =======================
func getPort() string {
	loadEnv()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
func getAddr() string {
	loadEnv()
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	addr := os.Getenv("ADDR")

	// 优先使用 ADDR
	if addr != "" {
		return addr
	}

	if port == "" {
		port = "8080"
	}
	if ip == "" {
		ip = "0.0.0.0"
	}
	return ip + ":" + port
}
func getenv(key, def string) string {
	loadEnv()
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	loadEnv()
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("环境变量 %s 格式错误: %v", key, err))
	}
	return i
}

func getEnvList(key string) []string {
	loadEnv()
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return []string{}
	}
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// ======================= 配置获取函数 =======================

func getEnv() string {
	loadEnv()
	env := os.Getenv("ENV")
	if env != LOCAL && env != PROD {
		panic(fmt.Sprintf("未知 ENV=%s", env))
	}
	return env
}

func getDBName() string {
	loadEnv()
	if ProdEnv() {
		dbName := os.Getenv("PROD_DB_NAME")
		if dbName == "" {
			panic("PROD_DB_NAME 未配置")
		}
		return dbName
	}
	dbName := os.Getenv("LOCAL_DB_NAME")
	if dbName == "" {
		panic("LOCAL_DB_NAME 未配置")
	}
	return dbName
}

func getDBUrl() string {
	loadEnv()
	if ProdEnv() {
		url := os.Getenv("PROD_DB_URL")
		if url == "" {
			panic("PROD_DB_URL 未配置")
		}
		return url
	}
	url := os.Getenv("LOCAL_DB_URL")
	if url == "" {
		panic("LOCAL_DB_URL 未配置")
	}
	return url
}

func getRedisAddr() string {
	loadEnv()
	if ProdEnv() {
		return os.Getenv("PROD_REDIS_ADDR")
	}
	return os.Getenv("LOCAL_REDIS_ADDR")
}

func getRedisPassword() string {
	loadEnv()
	if ProdEnv() {
		return os.Getenv("PROD_REDIS_PASSWORD")
	}
	return os.Getenv("LOCAL_REDIS_PASSWORD")
}

func getRedisDB() int {
	if ProdEnv() {
		return getEnvInt("PROD_REDIS_DB", 0)
	}
	return getEnvInt("LOCAL_REDIS_DB", 0)
}
func getAdminRExpirationDay() int64 {
	loadEnv()
	val := os.Getenv("ADMIN_R_EXPIRATION_DAY")
	if val == "" {
		log.Println("管理员 记住密码token 过期天数未配置, 使用默认值 30 天")
		val = "30"
	}
	d, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	log.Println(envLog+"管理员 记住密码token  过期天数:", d)
	return d
}
func getAdminExpirationDay() int64 {
	loadEnv()
	val := os.Getenv("ADMIN_EXPIRATION_DAY")
	if val == "" {
		log.Println("管理员 token 过期天数未配置, 使用默认值 5 天")
		val = "5"
	}
	d, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	log.Println(envLog+"管理员 token 过期天数:", d)
	return d
}

func getUserExpirationDay() int64 {
	loadEnv()
	val := os.Getenv("USER_EXPIRATION_DAY")
	if val == "" {
		log.Println("用户 token 过期天数未配置, 使用默认值 7 天")
		val = "7"
	}
	d, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	log.Println(envLog+"用户 token 过期天数:", d)
	return d
}

func getRootAdmins() map[string]struct{} {
	loadEnv()
	if ProdEnv() {
		return utils.ListToMap(getEnvList("PROD_ROOT_ADMINS"))
	}
	return utils.ListToMap(getEnvList("LOCAL_ROOT_ADMINS"))
}
