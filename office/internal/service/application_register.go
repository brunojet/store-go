package service

import (
	"context"

	"github.com/brunojet/store-go/office/internal/dto"
	"github.com/brunojet/store-go/shared/pkg/repo"
)

type ApplicationService struct {
	repo repo.ApplicationRepository
}

func NewApplicationService(r repo.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: r}
}

func (s *ApplicationService) Create(ctx context.Context, dto *dto.ApplicationCreate) error {
	return s.repo.Create(ctx, dto.ToDomain())
}

func (s *ApplicationService) GetByID(ctx context.Context, id int64) (*dto.Application, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return (*dto.Application)(c), nil
}

func (s *ApplicationService) Update(ctx context.Context, id int64, dto *dto.ApplicationUpdate) error {
	return s.repo.Update(ctx, dto.ToDomain((id)))
}

func (s *ApplicationService) List(ctx context.Context, p repo.ListParams) ([]dto.Application, error) {
	// Convert shared ListParams to infra-go ListParams and call ListWithParams
	ip := p.ToInfraParams()

	categories, _, err := s.repo.ListWithParams(ctx, ip, nil)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Application, len(categories))
	for i, c := range categories {
		result[i] = dto.Application(c)
	}
	return result, nil
}
