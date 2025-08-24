package repo

import (
	"context"

	internal "github.com/brunojet/store-go/infra/internal/repo"
	"gorm.io/gorm"
)

// Repository is a public wrapper around the infra internal repository implementation.
// We forward all calls to the internal implementation so callers can depend on this package
// instead of importing an internal package directly.
type Repository[T any] struct {
	inner *internal.Repository[T]
}

// NewRepository constructs a new public Repository that delegates to the internal implementation.
func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{inner: internal.NewRepository[T](db)}
}

func (r *Repository[T]) Create(ctx context.Context, ent *T) error { return r.inner.Create(ctx, ent) }
func (r *Repository[T]) GetByID(ctx context.Context, id int64, preloads ...string) (*T, error) {
	return r.inner.GetByID(ctx, id, preloads...)
}
func (r *Repository[T]) Update(ctx context.Context, ent *T) error { return r.inner.Update(ctx, ent) }
func (r *Repository[T]) Delete(ctx context.Context, ent *T) error { return r.inner.Delete(ctx, ent) }
func (r *Repository[T]) List(ctx context.Context, mod func(*gorm.DB) *gorm.DB) ([]T, error) {
	return r.inner.List(ctx, mod)
}
func (r *Repository[T]) ListWithParams(ctx context.Context, p *internal.ListParams, mod func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	return r.inner.ListWithParams(ctx, p, mod)
}
func (r *Repository[T]) FindWithParams(ctx context.Context, p *internal.ListParams, mod func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	return r.inner.FindWithParams(ctx, p, mod)
}
func (r *Repository[T]) WithTx(ctx context.Context, fn func(txRepo *Repository[T]) error) error {
	return r.inner.WithTx(ctx, func(txInner *internal.Repository[T]) error {
		return fn(&Repository[T]{inner: txInner})
	})
}
func (r *Repository[T]) DB() *gorm.DB { return r.inner.DB() }
