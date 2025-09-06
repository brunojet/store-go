package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ApplicationProfileHistory struct {
	BaseModel
	ApplicationContactId      int64                      `gorm:"column:application_contact_id"`
	ApplicationContact        ApplicationContact         `gorm:"foreignKey:ApplicationContactId"`
	ApplicationDetailId       int64                      `gorm:"column:application_detail_id"`
	ApplicationDetail         ApplicationDetail          `gorm:"foreignKey:ApplicationDetailId"`
	Categorys                 []Category                 `gorm:"many2many:ctgr_application_profile_history;joinForeignKey:id_application_profile_history;joinReferences:id_ctgr"`
	ReviewAt                  *time.Time                 `gorm:"column:review_at"`
	ProductionAt              *time.Time                 `gorm:"column:production_at"`
	DeactivatedAt             *time.Time                 `gorm:"column:deactivated_at"`
	DeactivationCause         *string                    `gorm:"column:deactivation_cause;size:255"`
	ApplicationConfigurations []ApplicationConfiguration `gorm:"-:migration;many2many:application_profile_history_configuration;joinForeignKey:ID;joinReferences:ApplicationId,IntegrationTypeId,TerminalModelId"`
}

func (ApplicationProfileHistory) TableName() string { return "application_profile_history" }

// O método PostMigrate é necessário devido à limitação do GORM em criar corretamente tabelas many-to-many
// com chaves primárias compostas. Aqui garantimos a criação manual da tabela de junção com PK e FKs corretas.
func (ApplicationProfileHistory) PostMigrate(db *gorm.DB) error {
	dialect := db.Dialector.Name()
	var idType string
	switch dialect {
	case "sqlite":
		idType = "INTEGER"
	case "postgres", "mysql":
		idType = "BIGINT"
	default:
		return fmt.Errorf("dialeto não suportado para criação manual da tabela many-to-many")
	}

	tableSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS application_profile_history_configuration (
		application_profile_history_id %s NOT NULL,
		application_id %s NOT NULL,
		integration_type_id %s NOT NULL,
		terminal_model_id %s NOT NULL,
		PRIMARY KEY (application_profile_history_id, application_id, integration_type_id, terminal_model_id),
		FOREIGN KEY (application_profile_history_id) REFERENCES application_profile_history(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
		FOREIGN KEY (application_id, integration_type_id, terminal_model_id) REFERENCES application_configuration(application_id, integration_type_id, terminal_model_id) ON DELETE RESTRICT ON UPDATE RESTRICT
		)`,
		idType, idType, idType, idType)

	if err := db.Exec(tableSQL).Error; err != nil {
		return fmt.Errorf("erro ao criar tabela de junção many-to-many: %w", err)
	}
	return nil
}
