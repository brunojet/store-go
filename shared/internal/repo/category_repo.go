package repo

import (
	"context"

	"github.com/brunojet/infra-go/pkg/repo"
	"github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

// Embed the generic repository to reuse Create/Update/GetByID/ListWithParams/DB
type CategoryRepo struct {
	repo.Repository[domain.Category]
}

func (r *CategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	return r.Repository.Create(ctx, c)
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	return r.Repository.GetByID(ctx, id)
}

func (r *CategoryRepo) Update(ctx context.Context, c *domain.Category) error {
	return r.Repository.Update(ctx, c)
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{Repository: repo.NewRepository[domain.Category](db)}
}

func (r *CategoryRepo) List(ctx context.Context, p ListParams) ([]domain.Category, error) {
	ip := p.ToInfraListParams()
	items, _, err := r.Repository.ListWithParams(ctx, ip, nil)
	return items, err
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	return r.Repository.Delete(ctx, &domain.Category{BaseEntity: domain.BaseEntity{BaseModel: domain.BaseModel{ID: id}}})
}
