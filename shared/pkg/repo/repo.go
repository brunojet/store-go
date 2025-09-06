package repo

import (
	"context"

	"gorm.io/gorm"

	internal "github.com/brunojet/store-go/shared/internal/repo"
	"github.com/brunojet/store-go/shared/pkg/domain"
)

type ListParams = internal.ListParams

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
	Update(ctx context.Context, c *domain.Category) error
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	List(ctx context.Context, p ListParams) ([]domain.Category, error)
	Delete(ctx context.Context, id int64) error
}

type CategoryTypeRepository interface {
	Create(ctx context.Context, t *domain.CategoryType) error
	Update(ctx context.Context, t *domain.CategoryType) error
	GetByID(ctx context.Context, id int64) (*domain.CategoryType, error)
	List(ctx context.Context, p ListParams) ([]domain.CategoryType, error)
	Delete(ctx context.Context, id int64) error
}

func NewCategoryRepo(db *gorm.DB) CategoryRepository {
	return internal.NewCategoryRepo(db)
}

func NewCategoryTypeRepo(db *gorm.DB) CategoryTypeRepository {
	return internal.NewCategoryTypeRepo(db)
}
