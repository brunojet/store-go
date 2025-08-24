package repo

import (
	"context"

	"store-go/shared/internal/domain"

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

// CreateWithImagem cria anexo e imagem em transação (Anexo deve ser criado primeiro quando se usa shared PK)
func (r *AnexoRepo) CreateWithImagem(ctx context.Context, anexo *domain.Anexo, imagem *domain.Imagem) error {
	// Usar transação para garantir atomicidade
	return r.WithTx(ctx, func(txRepo *Repository[domain.Anexo]) error {
		// criar anexo
		if err := txRepo.Create(ctx, anexo); err != nil {
			return err
		}
		// sincronizar ID da imagem com ID do anexo (shared PK)
		imagem.ID = anexo.ID
		if err := txRepo.DB().WithContext(ctx).Create(imagem).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetWithImagem carrega anexo com preload de imagem
func (r *AnexoRepo) GetWithImagem(ctx context.Context, id int64) (*domain.Anexo, error) {
	var a domain.Anexo
	if err := r.db.WithContext(ctx).Preload("Imagem").First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
