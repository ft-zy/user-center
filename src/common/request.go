package common

// 通用分页请求参数
type PageRequest struct {
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
}

// 通用删除请求
type DeleteRequest struct {
	Id int `json:"id"`
}
