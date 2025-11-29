package controller

import "github.com/StephenChristianW/go-movies-open/config"

// -------------------- 工具方法 --------------------

// PageSet 设置分页参数，保证合法性
func PageSet(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = config.DefaultPageSize
	}
	if pageSize < config.MinPageSize {
		pageSize = config.MinPageSize
	}
	if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}
	return page, pageSize
}
