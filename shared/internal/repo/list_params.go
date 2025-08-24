package repo

import "github.com/brunojet/infra-go/pkg/repo"

type ListParams struct {
	Page     int
	PageSize int
	OrderBy  string
	Desc     bool
	Filters  map[string]interface{}
}

func (p *ListParams) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
}

func (p *ListParams) ToInfraListParams() *repo.ListParams {
	p.Normalize()
	ip := repo.ListParams{
		Offset:      (p.Page - 1) * p.PageSize,
		Limit:       p.PageSize,
		OrderColumn: p.OrderBy,
		OrderDir:    "asc",
	}
	if p.Desc {
		ip.OrderDir = "desc"
	}
	return &ip
}
