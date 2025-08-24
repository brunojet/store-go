package repo

import (
	"context"

	"gorm.io/gorm"

	internal "github.com/brunojet/store-go/shared/internal/repo"
	"github.com/brunojet/store-go/shared/pkg/domain"
)

type ListParams = internal.ListParams

type CategoriaRepository interface {
	Create(ctx context.Context, c *domain.Categoria) error
	Update(ctx context.Context, c *domain.Categoria) error
	GetByID(ctx context.Context, id int64) (*domain.Categoria, error)
	List(ctx context.Context, p ListParams) ([]domain.Categoria, error)
	Delete(ctx context.Context, id int64) error
}

type TipoCategoriaRepository interface {
	Create(ctx context.Context, t *domain.TipoCategoria) error
	Update(ctx context.Context, t *domain.TipoCategoria) error
	GetByID(ctx context.Context, id int64) (*domain.TipoCategoria, error)
	List(ctx context.Context, p ListParams) ([]domain.TipoCategoria, error)
	Delete(ctx context.Context, id int64) error
}

func NewCategoriaRepo(db *gorm.DB) CategoriaRepository {
	return internal.NewCategoriaRepo(db)
}

func NewTipoCategoriaRepo(db *gorm.DB) TipoCategoriaRepository {
	return internal.NewTipoCategoriaRepo(db)
}
