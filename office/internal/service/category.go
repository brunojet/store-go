package service

import (
	"context"

	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
)

type CategoryService struct {
	repo repo.CategoryRepository
}

func NewCategoryService(r repo.CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) List(ctx context.Context, page, pageSize int) ([]domain.Category, error) {
	params := repo.ListParams{Page: page, PageSize: pageSize}
	return s.repo.List(ctx, params)
}
