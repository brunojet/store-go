package domain

type Stage int16

const (
	ApplicationCatalogTableName       = "application_catalog"
	StageDevelopment            Stage = 10
	StageTesting                Stage = 20
	StagePilot                  Stage = 30
	StageProduction             Stage = 40
)

type ApplicationCatalog struct {
	IntegrationTypeId    int64  `gorm:"primaryKey:pk_ctlg,priority:0;autoIncrement=false;index:ak_app_ctlg_u0,priority:2" json:"integration_type_id"`
	TerminalModelId      int64  `gorm:"primaryKey:pk_ctlg,priority:1;autoIncrement=false;index:ak_app_ctlg_u0,priority:3" json:"terminal_model_id"`
	Stage                Stage  `gorm:"primaryKey:pk_ctlg,priority:2;autoIncrement=false;index:ak_app_ctlg_u0,priority:1;check:stage IN (10,20,30,40)" json:"stage"`
	ApplicationId        int64  `gorm:"primaryKey:pk_ctlg,priority:3;autoIncrement=false;index:ak_app_ctlg_u0,priority:0" json:"application_id"`
	ApplicationProfileId *int64 `gorm:"not null;index:idx_ac_application_profile" json:"application_profile_id,omitempty"`
	ApplicationVersionId *int64 `gorm:"not null;index:idx_ac_application_version" json:"application_version_id,omitempty"`
	Active               *bool  `gorm:"index;default:false" json:"active,omitempty"`

	//Relacionamentos
	ApplicationProfile *ApplicationProfileHistory `gorm:"foreignKey:ApplicationProfileId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"application_profile,omitempty"`
	ApplicationVersion *ApplicationVersionHistory `gorm:"foreignKey:ApplicationVersionId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"application_version,omitempty"`
}

func (ApplicationCatalog) TableName() string { return ApplicationCatalogTableName }
