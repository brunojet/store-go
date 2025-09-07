package domain

import (
	"time"
)

const (
	ApplicationVersionTableName = "application_version"
)

type ApplicationVersionHistory struct {
	BaseEntity
	ApplicationId     int64      `gorm:"primaryKey;column:application_id"`
	IntegrationTypeId int64      `gorm:"primaryKey;column:integration_type_id"`
	TerminalModelId   int64      `gorm:"primaryKey;column:terminal_model_id"`
	Tamanho           int64      `gorm:"column:tamanho;not null"`
	NomeVersao        string     `gorm:"column:nome_versao;size:16;not null"`
	IdImage           int64      `gorm:"column:id_img;index"`
	PilotAt           *time.Time `gorm:"column:pilot_at;index"`
	ProductionAt      *time.Time `gorm:"column:production_at;index"`
	DeactivatedAt     *time.Time `gorm:"column:deactivated_at;index"`
	DeactivationCause *string    `gorm:"column:deactivation_cause;index;size:255"`

	//Relacionamentos
	Image               Image                `gorm:"foreignKey:IdImage;references:ID"`
	ApplicationCatalogs []ApplicationCatalog `gorm:"foreignKey:ApplicationVersionId,ApplicationId,IntegrationTypeId,TerminalModelId;references:ID,ApplicationId,IntegrationTypeId,TerminalModelId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationVersionHistory) TableName() string { return ApplicationVersionTableName }
