package repo

import (
	"context"

	internal "github.com/brunojet/store-go/infra/internal/repo"
	"github.com/brunojet/store-go/infra/pkg/domain"
	"gorm.io/gorm"
)

type ListParams = internal.ListParams

type Repository[T any] interface {
	Create(ctx context.Context, ent *T) error
	GetByID(ctx context.Context, id int64, preloads ...string) (*T, error)
	Update(ctx context.Context, ent *T) error
	Delete(ctx context.Context, ent *T) error
	ListWithParams(ctx context.Context, p *ListParams, mod func(*gorm.DB) *gorm.DB) ([]T, int64, error)
	DB() *gorm.DB
}

func NewRepository[T any](db *gorm.DB) Repository[T] {
	return internal.NewRepository[T](db)
}

type AnexoRepo interface {
	FindByNome(ctx context.Context, nome string) ([]domain.Anexo, error)
	CreateWith(ctx context.Context, anexo *domain.Anexo, createChild func(tx *gorm.DB, anexoID int64) error) error
	GetWith(ctx context.Context, id int64, preloads ...string) (*domain.Anexo, error)
	DB() *gorm.DB
}

func NewAnexoRepo(db *gorm.DB) AnexoRepo {
	return internal.NewAnexoRepo(db)
}
