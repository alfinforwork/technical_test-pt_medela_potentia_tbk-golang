package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type PaginationParams struct {
	Page     int
	PageSize int
	Search   string
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Data []interface{}  `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func GetPaginationParams(c fiber.Ctx) PaginationParams {
	page := 1
	pageSize := 10
	search := ""

	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if size, err := strconv.Atoi(ps); err == nil && size > 0 {
			if size > 100 {
				size = 100 // Max page size
			}
			pageSize = size
		}
	}

	search = c.Query("search")

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}
}

func CalculateTotalPages(total int64, pageSize int) int {
	if total == 0 {
		return 0
	}
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return pages
}
