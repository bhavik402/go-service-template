package util

import (
	"net/http"
	"strconv"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Page is the parsed pagination request: Limit/Offset for repositories,
// Page/PageSize for building response metadata.
type Page struct {
	Page     int
	PageSize int
}

func (p Page) Limit() int  { return p.PageSize }
func (p Page) Offset() int { return (p.Page - 1) * p.PageSize }

// ParsePage reads ?page= and ?page_size= from the request, applying sane
// defaults and bounds.
func ParsePage(r *http.Request) Page {
	page := parseIntDefault(r.URL.Query().Get("page"), 1)
	if page < 1 {
		page = 1
	}

	pageSize := parseIntDefault(r.URL.Query().Get("page_size"), DefaultPageSize)
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return Page{Page: page, PageSize: pageSize}
}

// PageMeta is the pagination metadata returned alongside list responses.
type PageMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func NewPageMeta(p Page, totalItems int) PageMeta {
	totalPages := totalItems / p.PageSize
	if totalItems%p.PageSize != 0 {
		totalPages++
	}
	return PageMeta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

func parseIntDefault(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}
