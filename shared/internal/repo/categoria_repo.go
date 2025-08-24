package repo

import (
	"context"

	"github.com/brunojet/store-go/infra/pkg/repo"
	"github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

type CategoriaRepo struct {
	rep repo.Repository[domain.Categoria]
}

func NewCategoriaRepo(db *gorm.DB) *CategoriaRepo {
	return &CategoriaRepo{rep: repo.NewRepository[domain.Categoria](db)}
}

func (r *CategoriaRepo) Create(ctx context.Context, c *domain.Categoria) error {
	return r.rep.Create(ctx, c)
}

func (r *CategoriaRepo) Update(ctx context.Context, c *domain.Categoria) error {
	return r.rep.Update(ctx, c)
}

func (r *CategoriaRepo) GetByID(ctx context.Context, id int64) (*domain.Categoria, error) {
	return r.rep.GetByID(ctx, id)
}

func (r *CategoriaRepo) List(ctx context.Context, p ListParams) ([]domain.Categoria, error) {
	ip := p.ToInfraListParams()
	items, _, err := r.rep.ListWithParams(ctx, ip, nil)
	return items, err
}

func (r *CategoriaRepo) Delete(ctx context.Context, id int64) error {
	return r.rep.DB().WithContext(ctx).Delete(&domain.Categoria{}, id).Error
}

func (r *CategoriaRepo) WithTx(tx interface{}) *CategoriaRepo {
	if gdb, ok := tx.(*gorm.DB); ok {
		return &CategoriaRepo{rep: repo.NewRepository[domain.Categoria](gdb)}
	}
	return r
}
