package services

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`       // 总条数
	TotalPages int `json:"total_pages"` // 总页数
	Data       any `json:"data"`
}
