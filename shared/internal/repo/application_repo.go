package repo

import (
	"context"

	"github.com/brunojet/infra-go/pkg/repo"
	"github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

// Embed the generic repository to reuse Create/Update/GetByID/ListWithParams/DB
type ApplicationRepo struct {
	repo.Repository[domain.Application]
}

func (r *ApplicationRepo) Create(ctx context.Context, c *domain.Application) error {
	return r.Repository.Create(ctx, c)
}

func (r *ApplicationRepo) GetByID(ctx context.Context, id int64) (*domain.Application, error) {
	return r.Repository.GetByID(ctx, id)
}

func (r *ApplicationRepo) Update(ctx context.Context, c *domain.Application) error {
	return r.Repository.Update(ctx, c)
}

func NewApplicationRepo(db *gorm.DB) *ApplicationRepo {
	return &ApplicationRepo{Repository: repo.NewRepository[domain.Application](db)}
}

func (r *ApplicationRepo) List(ctx context.Context, p ListParams) ([]domain.Application, error) {
	ip := p.ToInfraParams()
	items, _, err := r.Repository.ListWithParams(ctx, ip, nil)
	return items, err
}

func (r *ApplicationRepo) Delete(ctx context.Context, id int64) error {
	return r.Repository.Delete(ctx, &domain.Application{BaseEntity: domain.BaseEntity{BaseModel: domain.BaseModel{ID: id}}})
}
