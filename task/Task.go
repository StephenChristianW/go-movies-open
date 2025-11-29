package task

import (
	"log"
	"time"
)

const taskLog = "[TASK] "

// RunDaily 安全启动每日零点执行任务
// job: 要执行的函数，返回 error 表示是否成功
// taskName: 任务名称，用于日志
func RunDaily(taskName string, job func() error) {
	go func() {
		for {
			// 计算到下一天零点的时间间隔
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			sleepDuration := next.Sub(now)

			log.Printf(taskLog+"定时任务 [%s] 下一次执行时间: %s\n", taskName, next.Format("2006-01-02 15:04:05"))
			time.Sleep(sleepDuration)

			// 执行任务
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf(taskLog+"[%s] panic recovered: %v\n", taskName, r)
					}
				}()

				log.Printf(taskLog+"开始执行任务 [%s] 执行时间: %s", taskName, time.Now().Format("2006-01-02 15:04:05"))

				if err := job(); err != nil {
					log.Printf(taskLog+"[%s] 执行出错: %v\n", taskName, err)
				} else {
					log.Printf(taskLog+"[%s] 执行完成: %s\n", taskName, time.Now().Format("2006-01-02 15:04:05"))
				}
			}()
		}
	}()
}
