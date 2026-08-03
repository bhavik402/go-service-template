package util

import (
	"net/http"
	"net/url"
	"testing"
)

func TestParsePage_Defaults(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}

	page := ParsePage(req)

	if page.Page != 1 {
		t.Errorf("expected default page 1, got %d", page.Page)
	}
	if page.PageSize != DefaultPageSize {
		t.Errorf("expected default page size %d, got %d", DefaultPageSize, page.PageSize)
	}
}

func TestParsePage_ClampsPageSize(t *testing.T) {
	req := &http.Request{URL: &url.URL{RawQuery: "page=2&page_size=1000"}}

	page := ParsePage(req)

	if page.Page != 2 {
		t.Errorf("expected page 2, got %d", page.Page)
	}
	if page.PageSize != MaxPageSize {
		t.Errorf("expected page size clamped to %d, got %d", MaxPageSize, page.PageSize)
	}
	if got, want := page.Offset(), (2-1)*MaxPageSize; got != want {
		t.Errorf("expected offset %d, got %d", want, got)
	}
}

func TestNewPageMeta_RoundsUpTotalPages(t *testing.T) {
	page := Page{Page: 1, PageSize: 10}

	meta := NewPageMeta(page, 25)

	if meta.TotalPages != 3 {
		t.Errorf("expected 3 total pages for 25 items at page size 10, got %d", meta.TotalPages)
	}
}
