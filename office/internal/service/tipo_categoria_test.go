package service

import (
	"context"
	"testing"

	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
)

type mockTipoCategoriaRepo struct {
	receivedPage int
	receivedSize int
}

func (m *mockTipoCategoriaRepo) List(ctx context.Context, p repo.ListParams) ([]domain.TipoCategoria, error) {
	m.receivedPage = p.Page
	m.receivedSize = p.PageSize
	// return a simple slice
	return []domain.TipoCategoria{{}}, nil
}

func (m *mockTipoCategoriaRepo) Create(ctx context.Context, t *domain.TipoCategoria) error {
	return nil
}
func (m *mockTipoCategoriaRepo) Update(ctx context.Context, t *domain.TipoCategoria) error {
	return nil
}
func (m *mockTipoCategoriaRepo) GetByID(ctx context.Context, id int64) (*domain.TipoCategoria, error) {
	return &domain.TipoCategoria{}, nil
}
func (m *mockTipoCategoriaRepo) Delete(ctx context.Context, id int64) error { return nil }

func TestTipoCategoriaService_List_forwardsParamsAndReturnsItems(t *testing.T) {
	mock := &mockTipoCategoriaRepo{}
	svc := NewTipoCategoriaService(mock)

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
