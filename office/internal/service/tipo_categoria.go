package service

import (
	"context"

	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
)

// minimal repository interface used by the service (keeps dependency surface small)
type TipoCategoriaService struct {
	repo repo.TipoCategoriaRepository
}

func NewTipoCategoriaService(r repo.TipoCategoriaRepository) *TipoCategoriaService {
	return &TipoCategoriaService{repo: r}
}

func (s *TipoCategoriaService) List(ctx context.Context, page, pageSize int) ([]domain.TipoCategoria, error) {
	params := repo.ListParams{Page: page, PageSize: pageSize}
	return s.repo.List(ctx, params)
}
