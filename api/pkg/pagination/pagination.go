package pagination

import (
	"net/http"
	"strconv"
)

type Args struct {
	Page     int
	PageSize int
}

type Meta struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type Response struct {
	Data       interface{} `json:"data"`
	Pagination *Meta       `json:"pagination"`
}

func GetArgs(r *http.Request) *Args {
	page := 1
	pageSize := 50

	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}

	if v := r.URL.Query().Get("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return &Args{
		Page:     page,
		PageSize: pageSize,
	}
}

func NewResponse(data interface{}, total int, args *Args) *Response {
	totalPages := total / args.PageSize
	if total%args.PageSize != 0 {
		totalPages++
	}

	return &Response{
		Data: data,
		Pagination: &Meta{
			Page:       args.Page,
			PageSize:   args.PageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    args.Page < totalPages,
			HasPrev:    args.Page > 1,
		},
	}
}
