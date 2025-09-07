package domain

import (
	"time"
)

const (
	ApplicationVersionTableName = "application_version"
)

type ApplicationVersionHistory struct {
	BaseEntity
	ApplicationId     int64      `gorm:"primaryKey:pk_base,priority:0" json:"application_id"`
	IntegrationTypeId int64      `gorm:"primaryKey:pk_base,priority:1" json:"integration_type_id"`
	TerminalModelId   int64      `gorm:"primaryKey:pk_base,priority:2" json:"terminal_model_id"`
	VersionName       string     `gorm:"not null;size:100" json:"version_name"`
	VersionCode       int64      `gorm:"not null" json:"version_code"`
	Size              int64      `gorm:"not null" json:"size"`
	ImageId           int64      `gorm:"index" json:"image_id"`
	PilotAt           *time.Time `gorm:"index" json:"pilot_at,omitempty"`
	ProductionAt      *time.Time `gorm:"index" json:"production_at,omitempty"`
	DeactivatedAt     *time.Time `gorm:"index" json:"deactivated_at,omitempty"`
	DeactivationCause *string    `gorm:"index;size:255" json:"deactivation_cause,omitempty"`

	//Relacionamentos
	Image               *Image               `gorm:"foreignKey:ImageId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"image,omitempty"`
	ApplicationCatalogs []ApplicationCatalog `gorm:"foreignKey:ApplicationVersionId,ApplicationId,IntegrationTypeId,TerminalModelId;references:ID,ApplicationId,IntegrationTypeId,TerminalModelId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationVersionHistory) TableName() string { return ApplicationVersionTableName }
