package response

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// 分页结果
type PageGridData[T any] struct {
	Rows    []T `json:"rows"`
	Total   int `json:"total"`
	Summary any `json:"summary"` // 统计数据
}
