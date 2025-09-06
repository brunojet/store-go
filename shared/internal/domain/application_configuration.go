package domain

type ApplicationConfiguration struct {
	// Composite Primary Key (in this order): ApplicationId, IntegrationTypeId, TerminalModelId
	ApplicationId     int64           `gorm:"primaryKey,priority:1;column:application_id;not null;index:uk_cfg_rev,priority:3;uniqueIndex:uk_cfg_pk,priority:1"`
	Application       Application     `gorm:"foreignKey:ApplicationId;references:ID"`
	IntegrationTypeId int64           `gorm:"primaryKey,priority:2;column:integration_type_id;not null;index:uk_cfg_rev,priority:2;uniqueIndex:uk_cfg_pk,priority:2"`
	IntegrationType   IntegrationType `gorm:"foreignKey:IntegrationTypeId;references:ID"`
	TerminalModelId   int64           `gorm:"primaryKey,priority:3;column:terminal_model_id;not null;index:uk_cfg_rev,priority:1;uniqueIndex:uk_cfg_pk,priority:3"`
	TerminalModel     TerminalModel   `gorm:"foreignKey:TerminalModelId;references:ID"`

	// Relação reversa N:1 - Uma Configuration tem muitos Catalogs (por Stage)
	// -:migration evita que esta relação crie FK reversa; só ApplicationCatalog controla a FK
	ApplicationCatalogs []ApplicationCatalog        `gorm:"-:migration;foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId"`
	ApplicationVersions []ApplicationVersionHistory `gorm:"-:migration;foreignKey:ApplicationId,IntegrationTypeId,TerminalModelId;references:ApplicationId,IntegrationTypeId,TerminalModelId"`
}

func (ApplicationConfiguration) TableName() string { return "application_configuration" }
