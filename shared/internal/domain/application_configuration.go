package domain

type ApplicationConfiguration struct {
	// Composite Primary Key (in this order): ApplicationId, IntegrationTypeId, TerminalModelId
	ApplicationId     int64 `gorm:"primaryKey:pk_cfg,priority:1;index:uk_cfg_rev,priority:3;uniqueIndex:uk_cfg_pk,priority:1" json:"application_id"`
	IntegrationTypeId int64 `gorm:"primaryKey:pk_cfg,priority:2;index:uk_cfg_rev,priority:2;uniqueIndex:uk_cfg_pk,priority:2" json:"integration_type_id"`
	TerminalModelId   int64 `gorm:"primaryKey:pk_cfg,priority:3;index:uk_cfg_rev,priority:1;uniqueIndex:uk_cfg_pk,priority:3" json:"terminal_model_id"`

	Application         Application                 `gorm:"foreignKey:ApplicationId;references:ID" json:"application"`
	IntegrationType     IntegrationType             `gorm:"foreignKey:IntegrationTypeId;references:ID" json:"integration_type"`
	TerminalModel       TerminalModel               `gorm:"foreignKey:TerminalModelId;references:ID" json:"terminal_model"`
	ApplicationCatalogs []ApplicationCatalog        `gorm:"foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId" json:"application_catalogs"`
	ApplicationProfiles []ApplicationProfileHistory `gorm:"many2many:application_profile_history_configuration;joinForeignKey:ApplicationId,IntegrationTypeId,TerminalModelId;joinReferences:ProfileID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"application_profiles"`
	ApplicationVersions []ApplicationVersionHistory `gorm:"foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId" json:"application_versions"`
}

func (ApplicationConfiguration) TableName() string { return "application_configuration" }
