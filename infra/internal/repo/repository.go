package repo

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) withDB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *Repository[T]) applyPreloads(q *gorm.DB, preloads []string) *gorm.DB {
	for _, p := range preloads {
		q = q.Preload(p)
	}
	return q
}

func (r *Repository[T]) Create(ctx context.Context, ent *T) error {
	return r.db.WithContext(ctx).Create(ent).Error
}

func (r *Repository[T]) GetByID(ctx context.Context, id int64, preloads ...string) (*T, error) {
	var e T
	q := r.withDB(ctx)
	q = r.applyPreloads(q, preloads)
	if err := q.First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repository[T]) Update(ctx context.Context, ent *T) error {
	return r.db.WithContext(ctx).Save(ent).Error
}

func (r *Repository[T]) Delete(ctx context.Context, ent *T) error {
	return r.db.WithContext(ctx).Delete(ent).Error
}

func (r *Repository[T]) prepareModel(ctx context.Context, mod func(*gorm.DB) *gorm.DB) ([]T, *gorm.DB) {
	var out []T

	q := r.withDB(ctx).Model(new(T))

	if mod != nil {
		q = mod(q)
	}

	return out, q
}

func (r *Repository[T]) findList(out *[]T, tx *gorm.DB) error {
	if err := tx.Find(out).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository[T]) List(ctx context.Context, mod func(*gorm.DB) *gorm.DB) ([]T, error) {
	out, q := r.prepareModel(ctx, mod)
	err := r.findList(&out, q)
	return out, err
}

func (r *Repository[T]) ListWithParams(ctx context.Context, p *ListParams, mod func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	out, q := r.prepareModel(ctx, mod)

	countQ := q
	var total int64
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	} else if total == 0 {
		return make([]T, 0), 0, nil
	}

	q = r.applyPreloads(q, p.Preloads)
	q = applyListParams(q, p)

	err := r.findList(&out, q)

	return out, total, err
}

func (r *Repository[T]) FindWithParams(ctx context.Context, p *ListParams, mod func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	return r.ListWithParams(ctx, p, mod)
}

func (r *Repository[T]) WithTx(ctx context.Context, fn func(txRepo *Repository[T]) error) error {
	return r.withDB(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewRepository[T](tx))
	})
}

func (r *Repository[T]) DB() *gorm.DB { return r.db }
