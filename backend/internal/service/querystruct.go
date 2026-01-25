package service

import "gateway-service/internal/define"

type QueryRequest struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Keyword string `json:"keyword"`
}

func NewQueryRequest() *QueryRequest {
	return &QueryRequest{
		Page:    1,
		Size:    define.DefaultSize,
		Keyword: "",
	}
}
