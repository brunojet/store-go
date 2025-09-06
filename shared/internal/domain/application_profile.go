package domain

import (
	"strings"
	"time"

	"gorm.io/gorm"
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
	Categories                []Category                 `gorm:"-:migration;many2many:application_profile_history_category;joinForeignKey:ID;joinReferences:CategoryID"`
	ApplicationConfigurations []ApplicationConfiguration `gorm:"-:migration;many2many:application_profile_history_configuration;joinReferences:ID,ApplicationId,IntegrationTypeId,TerminalModelId;joinForeignKey:ID"`
}

func (ApplicationProfileHistory) TableName() string {
	return "application_profile_history"
}

// AssociateApplicationConfigurations faz associação em lote de várias ApplicationConfiguration
func (profile *ApplicationProfileHistory) AssociateApplicationConfigurations(db *gorm.DB, cfgs []*ApplicationConfiguration) error {
	if len(cfgs) == 0 {
		return nil
	}
	query := "INSERT INTO application_profile_history_configuration (id, application_id, integration_type_id, terminal_model_id) VALUES "
	vals := make([]interface{}, 0, len(cfgs)*4)
	placeholders := make([]string, 0, len(cfgs))
	for _, acfg := range cfgs {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		vals = append(vals, profile.ID, acfg.ApplicationId, acfg.IntegrationTypeId, acfg.TerminalModelId)
	}
	query += strings.Join(placeholders, ", ") + " ON CONFLICT DO NOTHING"
	return db.Exec(query, vals...).Error
}
