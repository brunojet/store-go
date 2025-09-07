package domain

import (
	"time"

	"gorm.io/gorm"

	domain_tools "github.com/brunojet/store-go/shared/internal/domain_utils"
)

type ApplicationProfileHistory struct {
	BaseModel
	ApplicationContactId int64      `gorm:"column:application_contact_id"`
	ApplicationDetailId  int64      `gorm:"column:application_detail_id"`
	ReviewAt             *time.Time `gorm:"column:review_at"`
	ProductionAt         *time.Time `gorm:"column:production_at"`
	DeactivatedAt        *time.Time `gorm:"column:deactivated_at"`
	DeactivationCause    *string    `gorm:"column:deactivation_cause;size:255"`

	//Relacionamentos
	ApplicationContact        ApplicationContact         `gorm:"foreignKey:ApplicationContactId"`
	ApplicationDetail         ApplicationDetail          `gorm:"foreignKey:ApplicationDetailId"`
	Categories                []Category                 `gorm:"many2many:application_profile_history_category;joinForeignKey:ProfileID;joinReferences:CategoryID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationConfigurations []ApplicationConfiguration `gorm:"many2many:application_profile_history_configuration;joinForeignKey:ProfileID;joinReferences:ApplicationId,IntegrationTypeId,TerminalModelId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationCatalogs       []ApplicationCatalog       `gorm:"foreignKey:ApplicationProfileId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationProfileHistory) TableName() string {
	return "application_profile_history"
}

func (ApplicationProfileHistory) PostMigrate(db *gorm.DB) error {
	if err := domain_tools.EnsureCascadeOnDelete(
		db,
		"application_profile_history_category",
		"profile_id",
		"application_profile_history",
		"id",
	); err != nil {
		return err
	}

	if err := domain_tools.EnsureCascadeOnDelete(
		db,
		"application_profile_history_configuration",
		"profile_id",
		"application_profile_history",
		"id",
	); err != nil {
		return err
	}

	return nil
}
