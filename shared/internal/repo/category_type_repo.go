package repo

import (
	"context"

	"github.com/brunojet/infra-go/pkg/repo"
	"github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

// Embed the generic repository to reuse Create/Update/GetByID/ListWithParams/DB
type CategoryTypeRepo struct {
	repo.Repository[domain.CategoryType]
}

func (r *CategoryTypeRepo) Create(ctx context.Context, t *domain.CategoryType) error {
	return r.Repository.Create(ctx, t)
}

func (r *CategoryTypeRepo) GetByID(ctx context.Context, id int64) (*domain.CategoryType, error) {
	return r.Repository.GetByID(ctx, id)
}

func (r *CategoryTypeRepo) Update(ctx context.Context, t *domain.CategoryType) error {
	return r.Repository.Update(ctx, t)
}

func NewCategoryTypeRepo(db *gorm.DB) *CategoryTypeRepo {
	return &CategoryTypeRepo{Repository: repo.NewRepository[domain.CategoryType](db)}
}

func (r *CategoryTypeRepo) List(ctx context.Context, p ListParams) ([]domain.CategoryType, error) {
	ip := p.ToInfraListParams()
	items, _, err := r.Repository.ListWithParams(ctx, ip, nil)
	return items, err
}

func (r *CategoryTypeRepo) Delete(ctx context.Context, id int64) error {
	return r.Repository.Delete(ctx, &domain.CategoryType{BaseEntity: domain.BaseEntity{BaseModel: domain.BaseModel{ID: id}}})
}
