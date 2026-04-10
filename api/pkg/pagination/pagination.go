package pagination

import (
	"net/http"
	"strconv"
)

type PaginationArgs struct {
	Page     int
	PageSize int
}

func GetPaginationArgs(r *http.Request) *PaginationArgs {
	userPage := r.URL.Query().Get("page")
	if userPage == "" {
		userPage = "1"
	}

	userPageSize := r.URL.Query().Get("page_size")
	if userPageSize == "" {
		userPageSize = "50"
	}

	page, err := strconv.Atoi(userPage)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(userPageSize)
	if err != nil || pageSize < 1 {
		pageSize = 50
	}

	return &PaginationArgs{
		Page:     page,
		PageSize: pageSize,
	}
}
