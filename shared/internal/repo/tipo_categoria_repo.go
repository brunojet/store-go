package repo

import (
	"context"

	"github.com/brunojet/infra-go/pkg/repo"
	"github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

type TipoCategoriaRepo struct {
	rep repo.Repository[domain.TipoCategoria]
}

func NewTipoCategoriaRepo(db *gorm.DB) *TipoCategoriaRepo {
	return &TipoCategoriaRepo{rep: repo.NewRepository[domain.TipoCategoria](db)}
}

func (r *TipoCategoriaRepo) Create(ctx context.Context, t *domain.TipoCategoria) error {
	return r.rep.Create(ctx, t)
}

func (r *TipoCategoriaRepo) Update(ctx context.Context, t *domain.TipoCategoria) error {
	return r.rep.Update(ctx, t)
}

func (r *TipoCategoriaRepo) GetByID(ctx context.Context, id int64) (*domain.TipoCategoria, error) {
	return r.rep.GetByID(ctx, id)
}

func (r *TipoCategoriaRepo) List(ctx context.Context, p ListParams) ([]domain.TipoCategoria, error) {
	ip := p.ToInfraListParams()
	items, _, err := r.rep.ListWithParams(ctx, ip, nil)
	return items, err
}

func (r *TipoCategoriaRepo) Delete(ctx context.Context, id int64) error {
	return r.rep.DB().WithContext(ctx).Delete(&domain.TipoCategoria{}, id).Error
}

func (r *TipoCategoriaRepo) WithTx(tx interface{}) *TipoCategoriaRepo {
	if gdb, ok := tx.(*gorm.DB); ok {
		return &TipoCategoriaRepo{rep: repo.NewRepository[domain.TipoCategoria](gdb)}
	}
	return r
}
