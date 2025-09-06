package domain

import (
	"fmt"

	domain_tools "github.com/brunojet/store-go/shared/internal/domain_utils"
	"gorm.io/gorm"
)

type Stage int16

const (
	ApplicationCatalogTableName       = "application_catalog"
	StageDevelopment            Stage = 0
	StageTesting                Stage = 10
	StagePilot                  Stage = 20
	StageProduction             Stage = 30
)

type ApplicationCatalog struct {
	IntegrationTypeId    int64  `gorm:"primaryKey,priority:0;index:ak_app_ctlg_u0,priority:3;not null;column:integration_type_id"`
	TerminalModelId      int64  `gorm:"primaryKey,priority:1;index:ak_app_ctlg_u0,priority:2;not null;column:terminal_model_id"`
	Stage                Stage  `gorm:"primaryKey,priority:2;index:ak_app_ctlg_u0,priority:1;not null;column:stage;check:stage IN (0,10,20,30)"`
	ApplicationId        int64  `gorm:"primaryKey,priority:3;index:ak_app_ctlg_u0,priority:0;not null;column:application_id"`
	ApplicationProfileId *int64 `gorm:"not null;column:application_profile_id"`
	ApplicationVersionId *int64 `gorm:"not null;column:application_version_id"`
	Ativo                *bool  `gorm:"index;default:false" json:"ativo"`

	// Relacionamentos
	ApplicationConfiguration ApplicationConfiguration  `gorm:"-:migration;foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId"`
	ApplicationProfile       ApplicationProfileHistory `gorm:"foreignKey:ApplicationProfileId;references:ID"`
	ApplicationVersion       ApplicationVersionHistory `gorm:"foreignKey:ApplicationVersionId;references:ID"`
}

func (ApplicationCatalog) TableName() string { return ApplicationCatalogTableName }

func (ApplicationCatalog) PostMigrate(db *gorm.DB) error {
	dialect := db.Dialector.Name()
	switch dialect {
	case "sqlite":
		if err := (ApplicationCatalog{}).recreateTable(db); err != nil {
			return err
		}
	case "postgres", "mysql":
		if err := (ApplicationCatalog{}).createPrimaryKey(db); err != nil {
			return err
		}
		if err := (ApplicationCatalog{}).createConstraints(db); err != nil {
			return err
		}
	default:
		return fmt.Errorf("dialeto não suportado para criação manual da constraint many-to-many")
	}
	return nil
}

func (ApplicationCatalog) recreateTable(db *gorm.DB) error {
	idType := "INTEGER"
	boolType := "BOOLEAN"

	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS application_catalog (
		integration_type_id %s NOT NULL,
		terminal_model_id %s NOT NULL,
		stage SMALLINT NOT NULL,
		application_id %s NOT NULL,
		application_profile_id %s NOT NULL,
		application_version_id %s NOT NULL,
		ativo %s DEFAULT FALSE
	);`,
		idType,
		idType,
		idType,
		idType,
		idType,
		boolType,
	)

	columns := []string{"integration_type_id", "terminal_model_id", "stage", "application_id", "application_profile_id", "application_version_id", "ativo"}

	return domain_tools.RecreateTable(db, ApplicationCatalogTableName, createTable, columns)
}

func (ApplicationCatalog) createPrimaryKey(db *gorm.DB) error {
	return domain_tools.CreatePrimaryKey(
		db,
		ApplicationCatalogTableName,
		"pk_app_catalog",
		[]string{"integration_type_id", "terminal_model_id", "stage", "application_id"},
	)
}

func (ApplicationCatalog) createConstraints(db *gorm.DB) error {
	return domain_tools.CreateConstraints(
		db,
		ApplicationCatalogTableName,
		"application_configuration",
		"fk_app_catalog_config",
		[]string{"application_id", "integration_type_id", "terminal_model_id"},
		"RESTRICT",
		"RESTRICT",
	)
}
