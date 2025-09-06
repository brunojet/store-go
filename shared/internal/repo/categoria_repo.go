package repo

import (
	"context"

	"github.com/brunojet/infra-go/pkg/repo"
	"github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	rep repo.Repository[domain.Category]
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{rep: repo.NewRepository[domain.Category](db)}
}

func (r *CategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	return r.rep.Create(ctx, c)
}

func (r *CategoryRepo) Update(ctx context.Context, c *domain.Category) error {
	return r.rep.Update(ctx, c)
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	return r.rep.GetByID(ctx, id)
}

func (r *CategoryRepo) List(ctx context.Context, p ListParams) ([]domain.Category, error) {
	ip := p.ToInfraListParams()
	items, _, err := r.rep.ListWithParams(ctx, ip, nil)
	return items, err
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	return r.rep.DB().WithContext(ctx).Delete(&domain.Category{}, id).Error
}

func (r *CategoryRepo) WithTx(tx interface{}) *CategoryRepo {
	if gdb, ok := tx.(*gorm.DB); ok {
		return &CategoryRepo{rep: repo.NewRepository[domain.Category](gdb)}
	}
	return r
}
