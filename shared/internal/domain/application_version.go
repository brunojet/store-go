package domain

import (
	"time"
)

const (
	ApplicationVersionTableName = "application_version"
)

type ApplicationVersionHistory struct {
	BaseEntity
	ApplicationId     int64      `gorm:"primaryKey:pk_base,priority:0"`
	IntegrationTypeId int64      `gorm:"primaryKey:pk_base,priority:1"`
	TerminalModelId   int64      `gorm:"primaryKey:pk_base,priority:2"`
	VersionName       string     `gorm:"not null;size:100"`
	VersionCode       int64      `gorm:"not null"`
	Size              int64      `gorm:"not null"`
	ImageId           int64      `gorm:"index"`
	PilotAt           *time.Time `gorm:"index"`
	ProductionAt      *time.Time `gorm:"index"`
	DeactivatedAt     *time.Time `gorm:"index"`
	DeactivationCause *string    `gorm:"index;size:255"`

	//Relacionamentos
	Image               Image                `gorm:"foreignKey:ImageId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationCatalogs []ApplicationCatalog `gorm:"foreignKey:ApplicationVersionId,ApplicationId,IntegrationTypeId,TerminalModelId;references:ID,ApplicationId,IntegrationTypeId,TerminalModelId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationVersionHistory) TableName() string { return ApplicationVersionTableName }
