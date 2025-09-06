package repo

import (
	"context"

	"github.com/brunojet/infra-go/pkg/repo"
	"github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

type CategoryTypeRepo struct {
	rep repo.Repository[domain.CategoryType]
}

func NewCategoryTypeRepo(db *gorm.DB) *CategoryTypeRepo {
	return &CategoryTypeRepo{rep: repo.NewRepository[domain.CategoryType](db)}
}

func (r *CategoryTypeRepo) Create(ctx context.Context, t *domain.CategoryType) error {
	return r.rep.Create(ctx, t)
}

func (r *CategoryTypeRepo) Update(ctx context.Context, t *domain.CategoryType) error {
	return r.rep.Update(ctx, t)
}

func (r *CategoryTypeRepo) GetByID(ctx context.Context, id int64) (*domain.CategoryType, error) {
	return r.rep.GetByID(ctx, id)
}

func (r *CategoryTypeRepo) List(ctx context.Context, p ListParams) ([]domain.CategoryType, error) {
	ip := p.ToInfraListParams()
	items, _, err := r.rep.ListWithParams(ctx, ip, nil)
	return items, err
}

func (r *CategoryTypeRepo) Delete(ctx context.Context, id int64) error {
	return r.rep.DB().WithContext(ctx).Delete(&domain.CategoryType{}, id).Error
}

func (r *CategoryTypeRepo) WithTx(tx interface{}) *CategoryTypeRepo {
	if gdb, ok := tx.(*gorm.DB); ok {
		return &CategoryTypeRepo{rep: repo.NewRepository[domain.CategoryType](gdb)}
	}
	return r
}
