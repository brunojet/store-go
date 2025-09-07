package domain

type Stage int16

const (
	ApplicationCatalogTableName       = "application_catalog"
	StageDevelopment            Stage = 0
	StageTesting                Stage = 10
	StagePilot                  Stage = 20
	StageProduction             Stage = 30
)

type ApplicationCatalog struct {
	IntegrationTypeId    int64  `gorm:"primaryKey:pk_ctlg,priority:0;autoIncrement=false;index:ak_app_ctlg_u0,priority:2"`
	TerminalModelId      int64  `gorm:"primaryKey:pk_ctlg,priority:1;autoIncrement=false;index:ak_app_ctlg_u0,priority:3"`
	Stage                Stage  `gorm:"primaryKey:pk_ctlg,priority:2;autoIncrement=false;index:ak_app_ctlg_u0,priority:1;check:stage IN (0,10,20,30)"`
	ApplicationId        int64  `gorm:"primaryKey:pk_ctlg,priority:3;autoIncrement=false;index:ak_app_ctlg_u0,priority:0"`
	ApplicationProfileId *int64 `gorm:"not null"`
	ApplicationVersionId *int64 `gorm:"not null"`
	Active               *bool  `gorm:"index;default:false"`

	//Relacionamentos
	ApplicationProfile ApplicationProfileHistory `gorm:"foreignKey:ApplicationProfileId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationVersion ApplicationVersionHistory `gorm:"foreignKey:ApplicationVersionId,ApplicationId,IntegrationTypeId,TerminalModelId;references:ID,ApplicationId,IntegrationTypeId,TerminalModelId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationCatalog) TableName() string { return ApplicationCatalogTableName }
