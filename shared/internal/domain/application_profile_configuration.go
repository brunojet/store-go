package domain

import (
	"gorm.io/gorm"

	domain_tools "github.com/brunojet/store-go/shared/internal/domain_utils"
)

type ApplicationProfileHistoryConfiguration struct {
	ID                int64 `gorm:"primaryKey:priority:0;column:id"`
	ApplicationId     int64 `gorm:"primaryKey:priority:1;column:application_id"`
	IntegrationTypeId int64 `gorm:"primaryKey:priority:2;column:integration_type_id"`
	TerminalModelId   int64 `gorm:"primaryKey:priority:3;column:terminal_model_id"`

	ApplicationProfileHistory ApplicationProfileHistory `gorm:"-:migration;foreignKey:ID;references:ID"`
	ApplicationConfiguration  ApplicationConfiguration  `gorm:"-:migration;foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId"`
}

func (ApplicationProfileHistoryConfiguration) TableName() string {
	return "application_profile_history_configuration"
}

// PostMigrate cria as constraints FK manualmente para a tabela de junção application_configuration
func (ApplicationProfileHistoryConfiguration) PostMigrate(db *gorm.DB) error {
	if err := domain_tools.CreateConstraints(
		db,
		"application_profile_history_configuration",
		"application_profile_history",
		"fk_app_profile_history_configuration_f0",
		[]string{"id"},
		"CASCADE",
		"RESTRICT",
	); err != nil {
		return err
	}

	if err := domain_tools.CreateConstraints(
		db,
		"application_profile_history_configuration",
		"application_configuration",
		"fk_app_profile_history_configuration_f1",
		[]string{"application_id", "integration_type_id", "terminal_model_id"},
		"RESTRICT",
		"RESTRICT",
	); err != nil {
		return err
	}

	return nil
}
