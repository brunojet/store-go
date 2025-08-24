package repo

import (
	"context"

	internal "github.com/brunojet/store-go/infra/internal/repo"
	"github.com/brunojet/store-go/infra/pkg/domain"
	"gorm.io/gorm"
)

// AnexoRepo is a public wrapper that forwards to the internal AnexoRepo implementation.
type AnexoRepo struct{ inner *internal.AnexoRepo }

func NewAnexoRepo(db *gorm.DB) *AnexoRepo { return &AnexoRepo{inner: internal.NewAnexoRepo(db)} }

func (r *AnexoRepo) FindByNome(ctx context.Context, nome string) ([]domain.Anexo, error) {
	return r.inner.FindByNome(ctx, nome)
}

func (r *AnexoRepo) CreateWith(ctx context.Context, anexo *domain.Anexo, createChild func(tx *gorm.DB, anexoID int64) error) error {
	return r.inner.CreateWith(ctx, anexo, createChild)
}

func (r *AnexoRepo) GetWith(ctx context.Context, id int64, preloads ...string) (*domain.Anexo, error) {
	return r.inner.GetWith(ctx, id, preloads...)
}
