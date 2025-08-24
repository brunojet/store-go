package repo

import (
	"context"

	"gorm.io/gorm"
)

// Repository é um repositório genérico baseado em GORM.
// T é o tipo da entidade (ex: domain.Anexo).
// Observação: o tipo T deve corresponder ao tipo usado nas chamadas (struct do domain).
type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// withDB centraliza r.db.WithContext(ctx) para facilitar testes e consistência.
func (r *Repository[T]) withDB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// applyPreloads aplica uma lista de preloads a uma query GORM.
// Centraliza o uso de Preload para facilitar testes reutilizáveis.
func (r *Repository[T]) applyPreloads(q *gorm.DB, preloads []string) *gorm.DB {
	for _, p := range preloads {
		q = q.Preload(p)
	}
	return q
}

// Create insere a entidade e retorna erro se houver.
func (r *Repository[T]) Create(ctx context.Context, ent *T) error {
	return r.db.WithContext(ctx).Create(ent).Error
}

// GetByID busca por chave primária (assume int64).
func (r *Repository[T]) GetByID(ctx context.Context, id int64, preloads ...string) (*T, error) {
	var e T
	q := r.withDB(ctx)
	q = r.applyPreloads(q, preloads)
	if err := q.First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

// Update salva (Save) a entidade inteira.
func (r *Repository[T]) Update(ctx context.Context, ent *T) error {
	return r.db.WithContext(ctx).Save(ent).Error
}

// Delete remove a entidade.
func (r *Repository[T]) Delete(ctx context.Context, ent *T) error {
	return r.db.WithContext(ctx).Delete(ent).Error
}

// List executa uma query customizável via modifier (preloads/where/order/limit).
func (r *Repository[T]) List(ctx context.Context, mod func(*gorm.DB) *gorm.DB) ([]T, error) {
	var out []T
	q := r.withDB(ctx).Model(new(T))
	if mod != nil {
		q = mod(q)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

// ListParams encapsula opções comuns para listagens: ordenação, paginação e preloads.
// Filtros e condições complexas são passados via callback `mod`.

// ApplyListParams aplica ordenação, preloads, offset e limit a uma query GORM.
// Condições WHERE devem ser aplicadas no callback `mod` passado a ListWithParams.
// ... ApplyListParams and buildOrder moved to repository_params.go

// ListWithParams retorna lista paginada baseada em ListParams e o total (antes de offset/limit).
// Recebe um modificador opcional `mod` para permitir joins/condições complexas antes de aplicar filtros/paginacao.
func (r *Repository[T]) ListWithParams(ctx context.Context, p *ListParams, mod func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	var out []T
	base := r.withDB(ctx).Model(new(T))
	// aplicar modificador (joins, selects customizados) antes de contar e buscar
	if mod != nil {
		base = mod(base)
	}
	// contar total; o modificador `mod` já foi aplicado a `base` e deve conter
	// quaisquer condições/joins necessários (WHERE).
	countQ := base
	var total int64
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := base
	// aplicar preloads/ordenacao/offset/limit
	// ApplyListParams aplica defaults e caps para limit/offset
	// aplicar preloads/ordem/offset/limit (ApplyListParams cuida de defaults e caps)
	q = r.applyPreloads(q, p.Preloads)
	q = applyListParams(q, p)
	if err := q.Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// FindWithParams é um wrapper conveniente que delega para ListWithParams.
// Mantém a mesma assinatura, apenas nome mais curto para uso comum em repositórios.
func (r *Repository[T]) FindWithParams(ctx context.Context, p *ListParams, mod func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	return r.ListWithParams(ctx, p, mod)
}

// WithTx executa uma função dentro de transação e provê um repo que usa a tx.
func (r *Repository[T]) WithTx(ctx context.Context, fn func(txRepo *Repository[T]) error) error {
	return r.withDB(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewRepository[T](tx))
	})
}

// DB expõe o gorm.DB subjacente quando necessário para queries avançadas.
func (r *Repository[T]) DB() *gorm.DB { return r.db }
