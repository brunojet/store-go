package repo

import (
	"context"

	"github.com/brunojet/store-go/infra/pkg/domain"
	"gorm.io/gorm"
)

// AnexoRepo é um wrapper com queries específicas para Anexo
type AnexoRepo struct {
	*Repository[domain.Anexo]
}

func NewAnexoRepo(db *gorm.DB) *AnexoRepo {
	return &AnexoRepo{Repository: NewRepository[domain.Anexo](db)}
}

// FindByNome retorna anexos com o nome informado
func (r *AnexoRepo) FindByNome(ctx context.Context, nome string) ([]domain.Anexo, error) {
	mod := func(q *gorm.DB) *gorm.DB {
		return q.Where("nome = ?", nome)
	}
	return r.List(ctx, mod)
}

// CreateWith cria um Anexo e delega a criação de um resource filho para um callback.
// O callback é executado dentro da mesma transação e recebe o *gorm.DB da transação
// e o ID do anexo recém-criado. Mantém AnexoRepo genérico e sem dependência de tipos concretos.
func (r *AnexoRepo) CreateWith(ctx context.Context, anexo *domain.Anexo, createChild func(tx *gorm.DB, anexoID int64) error) error {
	// Usar transação para garantir atomicidade
	return r.WithTx(ctx, func(txRepo *Repository[domain.Anexo]) error {
		// criar anexo
		if err := txRepo.Create(ctx, anexo); err != nil {
			return err
		}
		// delegar criação do child ao callback, usando a DB da transação
		if createChild != nil {
			if err := createChild(txRepo.DB().WithContext(ctx), anexo.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetWith carrega um Anexo com preloads opcionais.
func (r *AnexoRepo) GetWith(ctx context.Context, id int64, preloads ...string) (*domain.Anexo, error) {
	return r.GetByID(ctx, id, preloads...)
}
