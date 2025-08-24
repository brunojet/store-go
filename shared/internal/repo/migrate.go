package repo

import (
	"log"

	"gorm.io/gorm"

	"store-go/shared/internal/domain"
)

// MigrateAll aplica AutoMigrate na ordem segura para evitar constraints invertidas.
// Ordem sugerida: Anexo -> Imagem -> VersaoAplicativo -> demais tabelas que dependem delas.
func MigrateAll(db *gorm.DB) error {
	log.Println("running AutoMigrate: Anexo, Imagem, VersaoAplicativo...")
	// chamar AutoMigrate explicitamente por tipo (evita ambiguidade)
	if err := db.AutoMigrate(&domain.Anexo{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&domain.Imagem{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&domain.VersaoAplicativo{}); err != nil {
		return err
	}
	log.Println("migrations completed")
	return nil
}
