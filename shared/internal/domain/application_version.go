package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	domain_tools "github.com/brunojet/store-go/shared/internal/domain_utils"
)

const (
	ApplicationVersionTableName = "application_version"
)

type ApplicationVersion struct {
	BaseEntity
	ApplicationId            int64                    `gorm:"column:application_id;index:ak_app_ver_u0,priority:0"`
	IntegrationTypeId        int64                    `gorm:"column:integration_type_id;index:ak_app_ver_u0,priority:1"`
	TerminalModelId          int64                    `gorm:"column:terminal_model_id;index:ak_app_ver_u0,priority:2"`
	PilotAt                  *time.Time               `gorm:"column:pilot_at;index"`
	ProductionAt             *time.Time               `gorm:"column:production_at;index"`
	DeactivatedAt            *time.Time               `gorm:"column:deactivated_at;index"`
	DeactivationCause        *string                  `gorm:"column:deactivation_cause;index;size:255"`
	Tamanho                  int64                    `gorm:"column:tamanho;not null"`
	NomeVersao               string                   `gorm:"column:nome_versao;size:16;not null"`
	IdImage                  int64                    `gorm:"column:id_img;index"`
	Image                    Image                    `gorm:"foreignKey:IdImage;references:ID"`
	ApplicationConfiguration ApplicationConfiguration `gorm:"-:migration;foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId"`
}

func (ApplicationVersion) TableName() string { return ApplicationVersionTableName }

// O método PostMigrate é necessário para garantir PK e criar a FK composta manualmente devido à limitação do GORM.
func (ApplicationVersion) PostMigrate(db *gorm.DB) error {
	dialect := db.Dialector.Name()
	switch dialect {
	// Added missing opening curly brace for switch block
	case "sqlite":
		if err := (ApplicationVersion{}).recreateTableSql(db); err != nil {
			return err
		}
	case "postgres", "mysql":
		if err := (ApplicationVersion{}).createConstraints(db); err != nil {
			return err
		}
	default:
		return fmt.Errorf("dialeto não suportado para criação manual da constraint many-to-many")
	}
	return nil
}

func (ApplicationVersion) createConstraints(db *gorm.DB) error {
	return domain_tools.CreateConstraints(
		db,
		ApplicationVersionTableName,
		"application_configuration",
		"fk_application_configuration",
		[]string{"application_id", "integration_type_id", "terminal_model_id"},
	)
}

func (ApplicationVersion) recreateTableSql(db *gorm.DB) error {
	dialect := db.Dialector.Name()
	var idType, autoInc, timeType string
	if dialect != "sqlite" {
		return fmt.Errorf("dialeto não suportado para recriar a tabela application_version")
	}
	idType = "INTEGER"
	autoInc = "PRIMARY KEY AUTOINCREMENT"
	timeType = "DATETIME"

	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id %s,
		created_at %s NOT NULL,
		updated_at %s NOT NULL,
		nome VARCHAR(255) NOT NULL,
		descricao TEXT,
		ativo BOOLEAN NOT NULL DEFAULT FALSE,
		application_id %s NOT NULL,
		integration_type_id %s NOT NULL,
		terminal_model_id %s NOT NULL,
		pilot_at %s,
		production_at %s,
		deactivated_at %s,
		deactivation_cause VARCHAR(255),
		tamanho BIGINT NOT NULL,
		nome_versao VARCHAR(16) NOT NULL,
		id_img %s,
		FOREIGN KEY (id_img) REFERENCES image(id) ON DELETE SET NULL ON UPDATE CASCADE
	)`,
		ApplicationVersionTableName,
		autoInc,
		timeType, timeType,
		idType, idType, idType,
		timeType, timeType, timeType,
		idType,
	)

	columns := []string{"id", "created_at", "updated_at", "nome", "descricao", "ativo", "application_id", "integration_type_id", "terminal_model_id", "pilot_at", "production_at", "deactivated_at", "deactivation_cause", "tamanho", "nome_versao", "id_img"}

	if err := domain_tools.RecreateTable(db, ApplicationVersionTableName, createTable, columns); err != nil {
		return err
	}

	return nil
}
