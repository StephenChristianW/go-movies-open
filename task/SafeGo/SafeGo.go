package SafeGo

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

// LogPanic 将 panic 信息写入 exe 当前目录下的日志文件
func LogPanic(r interface{}) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("panic_%s.log", timestamp)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法创建 panic 文件:", err)
		return
	}
	defer f.Close()

	f.WriteString("========== PANIC ==========\n")
	f.WriteString(fmt.Sprintf("Panic: %v\n", r))
	f.WriteString(string(debug.Stack()))
	f.WriteString("\n===========================\n")

	fmt.Println("程序发生 panic，已记录到", filename)
}

// SafeGo 包装 goroutine，自动捕获 panic
func SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				LogPanic(r)
			}
		}()
		fn()
	}()
}

// RunWithRecovery 用于主函数中运行任意函数，捕获 panic 并阻塞等待
func RunWithRecovery(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			LogPanic(r)
			fmt.Println("程序发生 panic，按回车退出...")
			fmt.Scanln()
		}
	}()
	fn()
}
