package domain

import (
	"strings"

	id "github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

// MigrateApplicationConfigWithCatalogFK ensures proper migration order and FK creation
// between ApplicationConfiguration (parent) and ApplicationCatalog (child).
func MigrateApplicationConfigWithCatalogFK(db *gorm.DB) error {
	// Step 1: Migrate parent table first (ApplicationConfiguration) SEM FK automática
	if err := db.AutoMigrate(&id.ApplicationConfiguration{}); err != nil {
		return err
	}

	// Step 2: Migrate child table (ApplicationCatalog) SEM FK automática
	if err := db.AutoMigrate(&id.ApplicationCatalog{}); err != nil {
		return err
	}

	// Step 3: Criar FK manualmente APENAS de ApplicationCatalog -> ApplicationConfiguration
	constraintName := "fk_application_catalog_configuration"

	// Tentar criar constraint - usar GORM Migrator para melhor controle
	if !db.Migrator().HasConstraint(&id.ApplicationCatalog{}, constraintName) {
		// Criar constraint na direção correta: ApplicationCatalog -> ApplicationConfiguration
		sql := `ALTER TABLE application_catalog
				ADD CONSTRAINT ` + constraintName + `
				FOREIGN KEY (application_id, integration_type_id, terminal_model_id)
				REFERENCES application_configuration(application_id, integration_type_id, terminal_model_id)
				ON DELETE RESTRICT ON UPDATE RESTRICT`

		if err := db.Exec(sql).Error; err != nil {
			// Se falhar, pode ser erro de constraint já existir ou outro problema
			if !strings.Contains(strings.ToLower(err.Error()), "already exists") && !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				return err
			}
		}
	}

	return nil
}
