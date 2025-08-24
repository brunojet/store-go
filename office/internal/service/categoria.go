package service

import (
	"context"

	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
)

type CategoriaService struct {
	repo repo.CategoriaRepository
}

func NewCategoriaService(r repo.CategoriaRepository) *CategoriaService {
	return &CategoriaService{repo: r}
}

func (s *CategoriaService) List(ctx context.Context, page, pageSize int) ([]domain.Categoria, error) {
	params := repo.ListParams{Page: page, PageSize: pageSize}
	return s.repo.List(ctx, params)
}
