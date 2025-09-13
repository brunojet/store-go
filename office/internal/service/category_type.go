package service

import (
	"context"

	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
)

// minimal repository interface used by the service (keeps dependency surface small)
type CategoryTypeService struct {
	repo repo.CategoryTypeRepository
}

func NewCategoryTypeService(r repo.CategoryTypeRepository) *CategoryTypeService {
	return &CategoryTypeService{repo: r}
}

func (s *CategoryTypeService) List(ctx context.Context, page, pageSize int) ([]domain.CategoryType, error) {
	params := repo.ListParams{Page: page, PageSize: pageSize}
	return s.repo.List(ctx, params)
}
