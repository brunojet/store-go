package repo

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	DefaultLimit       = 50
	MaxLimit           = 1000
	DefaultOrderColumn = "created_at"
	DefaultOrderDir    = "desc"
)

var allowedColumns = map[string]bool{
	"id":         true,
	"name":       true,
	"nome":       true,
	"created_at": true,
}

type ListParams struct {
	Order       string
	Offset      int
	Limit       int
	Preloads    []string
	OrderColumn string
	OrderDir    string
}

func applyListParams(q *gorm.DB, p *ListParams) *gorm.DB {
	if col, desc := buildOrder(p); col != "" {
		q = q.Order(clause.OrderByColumn{Column: clause.Column{Name: col}, Desc: desc})
	}
	limit, offset := buildPageable(p)
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	return q
}

func buildPageable(p *ListParams) (limit, offset int) {
	if p == nil {
		return DefaultLimit, 0
	}
	limit = p.Limit
	if limit <= 0 {
		limit = DefaultLimit
	} else if limit > MaxLimit {
		limit = MaxLimit
	}
	offset = p.Offset
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func buildOrder(p *ListParams) (string, bool) {
	// handle nil p by using defaults
	if p == nil {
		return DefaultOrderColumn, strings.ToLower(DefaultOrderDir) == "desc"
	}
	col := p.OrderColumn
	dir := p.OrderDir
	if col == "" && p.Order != "" {
		parts := strings.Fields(p.Order)
		if len(parts) >= 1 {
			col = parts[0]
		}
		if len(parts) >= 2 {
			dir = parts[len(parts)-1]
		}
	}
	if col == "" {
		return "", false
	}
	if !allowedColumns[col] {
		return "", false
	}
	if dir == "" {
		dir = "asc"
	}
	dir = strings.ToLower(dir)
	if dir != "asc" && dir != "desc" {
		return "", false
	}
	return col, dir == "desc"
}
