package dto

type PageIn interface{}

// PageInput 通用分页输入结构
type PageInput struct {
	//每页条数
	PageSize int `json:"pageSize"`
	//当前页码
	PageNo int `json:"pageNo"`
}
