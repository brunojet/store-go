package service

import (
	"context"
	"testing"

	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
)

type mockCategoryTypeRepo struct {
	receivedPage int
	receivedSize int
}

func (m *mockCategoryTypeRepo) List(ctx context.Context, p repo.ListParams) ([]domain.CategoryType, error) {
	m.receivedPage = p.Page
	m.receivedSize = p.PageSize
	// return a simple slice
	return []domain.CategoryType{{}}, nil
}

func (m *mockCategoryTypeRepo) Create(ctx context.Context, t *domain.CategoryType) error {
	return nil
}
func (m *mockCategoryTypeRepo) Update(ctx context.Context, t *domain.CategoryType) error {
	return nil
}
func (m *mockCategoryTypeRepo) GetByID(ctx context.Context, id int64) (*domain.CategoryType, error) {
	return &domain.CategoryType{}, nil
}
func (m *mockCategoryTypeRepo) Delete(ctx context.Context, id int64) error { return nil }

func TestCategoryTypeService_List_forwardsParamsAndReturnsItems(t *testing.T) {
	mock := &mockCategoryTypeRepo{}
	svc := NewCategoryTypeService(mock)

	items, err := svc.List(context.Background(), 2, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if mock.receivedPage != 2 || mock.receivedSize != 50 {
		t.Fatalf("expected page=2 size=50, got page=%d size=%d", mock.receivedPage, mock.receivedSize)
	}
}
