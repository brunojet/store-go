package domain

type ApplicationConfiguration struct {
	// Composite Primary Key (in this order): ApplicationId, IntegrationTypeId, TerminalModelId
	ApplicationId     int64 `gorm:"primaryKey:pk_cfg,priority:1;index:uk_cfg_rev,priority:3;uniqueIndex:uk_cfg_pk,priority:1;column:application_id;not null"`
	IntegrationTypeId int64 `gorm:"primaryKey:pk_cfg,priority:2;index:uk_cfg_rev,priority:2;uniqueIndex:uk_cfg_pk,priority:2;column:integration_type_id;not null"`
	TerminalModelId   int64 `gorm:"primaryKey:pk_cfg,priority:3;index:uk_cfg_rev,priority:1;uniqueIndex:uk_cfg_pk,priority:3;column:terminal_model_id;not null"`

	Application         Application                 `gorm:"foreignKey:ApplicationId;references:ID"`
	IntegrationType     IntegrationType             `gorm:"foreignKey:IntegrationTypeId;references:ID"`
	TerminalModel       TerminalModel               `gorm:"foreignKey:TerminalModelId;references:ID"`
	ApplicationCatalogs []ApplicationCatalog        `gorm:"foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId"`
	ApplicationProfiles []ApplicationProfileHistory `gorm:"many2many:application_profile_history_configuration;joinForeignKey:ApplicationId,IntegrationTypeId,TerminalModelId;joinReferences:ProfileID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationVersions []ApplicationVersionHistory `gorm:"foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId"`
}

func (ApplicationConfiguration) TableName() string { return "application_configuration" }
