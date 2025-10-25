package dto

import "github.com/brunojet/store-go/shared/pkg/domain"

type ApplicationCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a *ApplicationCreate) ToDomain() *domain.Application {
	return &domain.Application{
		BaseEntity: domain.BaseEntity{
			Name:        a.Name,
			Description: a.Description,
		},
	}
}

type ApplicationUpdate struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (a *ApplicationUpdate) ToDomain(id int64) *domain.Application {
	return &domain.Application{
		BaseEntity: domain.BaseEntity{
			BaseModel:   domain.BaseModel{ID: id},
			Name:        *a.Name,
			Description: *a.Description,
		},
	}
}

type Application domain.Application
