package UtilsTime

import (
	"fmt"
	"time"
)

const format = "2006-01-02 15:04:05"

// NowTime 返回当前时间，格式化为 "YYYY-MM-DD HH:MM:SS"
func NowTime() string {
	return time.Now().Format(format)
}

// AfterNDays000 返回当前时间往后 after 天的零点时间，格式化为 "YYYY-MM-DD 00:00:00"
// after: 往后推的天数，可以为负数表示往前
func AfterNDays000(after int) string {
	t := time.Now().AddDate(0, 0, after)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Format(format)
}

// Today000 返回今天的零点时间，格式化为 "YYYY-MM-DD 00:00:00"
func Today000() string {
	t := time.Now().AddDate(0, 0, 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Format(format)
}

// DaysToTimeDuration 将天数转换为对应的 time.Duration
func DaysToTimeDuration(days int64) time.Duration {
	return time.Hour * 24 * time.Duration(days)
}
func PerDaySeconds(days int64) int64 {
	return int64(24 * 60 * 60 * days)
}

// CheckUnlockTime 判断是否已解锁
// 返回 (已解锁, 剩余秒数)
func CheckUnlockTime(t time.Time) (bool, int) {
	now := time.Now()
	if !t.After(now) {
		return true, 0
	}
	return false, int(t.Sub(now).Seconds()) // 未解锁，返回剩余秒数
}

// ParseFlexibleTime 尝试将多种常见格式的时间字符串解析为 time.Time。
// 支持的格式：
// - 2006-01-02
// - 2006-01-02 15:04:05
// - 2006/01/02 15:04:05
// - RFC3339 (如 2025-10-31T12:00:00Z)
//
// 返回值：
// - time.Time：解析成功的时间对象（失败则返回零值）
// - error：解析失败时返回错误信息
func ParseFlexibleTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"2006-01-02",
		"2006/01/02",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析时间格式: %s", value)
}
