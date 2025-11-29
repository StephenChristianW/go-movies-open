package initRootAdmin

import (
	"context"
	"errors"
	"fmt"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/StephenChristianW/go-movies-open/db/MongoDBServices/collections"
	"github.com/StephenChristianW/go-movies-open/services/System/Admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func InitRootAdmins() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var (
		notCreated = make(map[string]string)
		needOutput = false
	)

	// 检查是否需要输出（有需要创建的账户或错误）
	for name := range config.RootAdmins {
		if len(name) < config.LessLenAdminName {
			notCreated[name] = fmt.Sprintf("用户名至少%d个字符", config.LessLenAdminName)
			needOutput = true
			continue
		}

		if !adminExists(ctx, name) {
			needOutput = true
		}
	}

	// 只有需要时才输出
	if needOutput {
		fmt.Println("==============================================")
		fmt.Println("          开始初始化根管理员                ")
		fmt.Println("==============================================")
		fmt.Println()
		time.Sleep(1 * time.Second)
	}

	// 实际处理逻辑
	for name := range config.RootAdmins {
		if len(name) < config.LessLenAdminName {
			if needOutput {
				fmt.Printf("[ERROR] %s : 用户名至少%d个字符\n", name, config.LessLenAdminName)
				time.Sleep(300 * time.Millisecond)
			}
			continue
		}

		if !adminExists(ctx, name) {
			if needOutput {
				fmt.Printf("[CREATE] %s", name)
			}
			createAdminWithPassword(name)
			if needOutput {
				fmt.Println(" : DONE")
				time.Sleep(500 * time.Millisecond)
			}
		} else if needOutput {
			fmt.Printf("[EXISTS] %s : 跳过创建\n", name)
			time.Sleep(200 * time.Millisecond)
		}
	}

	// 错误汇总
	if len(notCreated) > 0 && needOutput {
		time.Sleep(1 * time.Second)
		fmt.Println()
		fmt.Println("------------ 创建失败汇总 ------------")
		for name, reason := range notCreated {
			fmt.Printf("[FAILED] %s : %s\n", name, reason)
			time.Sleep(500 * time.Millisecond)
		}
	}

	// 最终检查
	count, err := collections.GetAdminCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		if needOutput {
			fmt.Printf("\n[FATAL] %s\n", err.Error())
			time.Sleep(1 * time.Second)
		}
		panic(err)
	}

	if count == 0 {
		if needOutput {
			fmt.Printf("\n[FATAL] 没有成功创建任何根管理员\n")
			time.Sleep(1 * time.Second)
		}
		panic("错误: 没有成功创建任何根管理员")
	}

	// 只有有输出时才显示结束信息
	if needOutput {
		time.Sleep(1 * time.Second)
		fmt.Println()
		fmt.Println("==============================================")
		fmt.Println("          根管理员初始化完成                ")
		fmt.Println("==============================================")
		fmt.Println()
		time.Sleep(1 * time.Second)
	}
}
func adminExists(ctx context.Context, name string) bool {
	result := collections.GetAdminCollection().FindOne(ctx, bson.M{"username": name})
	if result.Err() == nil {
		return true
	}

	if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
		panic(result.Err())
	}
	return false
}

func createAdminWithPassword(name string) {
	fmt.Printf("\n正在创建初始管理员: %s\n", name)
	for {
		pwd := promptPassword(name)
		if isValidPassword(pwd) {
			createAdmins(name, pwd)
			fmt.Printf("管理员 %s 创建成功!\n", name)
			return
		}
		fmt.Printf("密码必须至少%d个字符，请重新输入\n", config.LessLenPwd)
	}
}

func promptPassword(name string) string {
	fmt.Printf("请输入 %s 的密码: ", name)
	var pwd string
	_, _ = fmt.Scanf("%s\n", &pwd)
	return pwd
}

func isValidPassword(pwd string) bool {
	return pwd != "" && len(pwd) >= config.LessLenPwd
}
func createAdmins(username string, password string) {
	var adminService AdminService.AdminInterface = &AdminService.AdminService{}
	err := adminService.CreateAdmin(username, password, "")
	if err != nil {
		panic(err)
	}
}
