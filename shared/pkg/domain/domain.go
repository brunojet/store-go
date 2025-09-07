package domain

// This package re-exports the internal domain types as public aliases so
// other modules (for example `office`) can depend on the domain without
// importing internal packages directly.

import (
	infra "github.com/brunojet/infra-go/pkg/domain"
	id "github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/gorm"
)

// Exportação das principais entidades para uso externo
type BaseModel = infra.BaseModel
type BaseEntity = infra.BaseEntity
type StorageObject = id.StorageObject
type Image = id.Image
type Video = id.Video
type IntegrationType = id.IntegrationType
type TerminalModel = id.TerminalModel
type Application = id.Application
type CategoryType = id.CategoryType
type Category = id.Category
type ApplicationContact = id.ApplicationContact
type ApplicationDetail = id.ApplicationDetail
type ApplicationVersionHistory = id.ApplicationVersionHistory
type ApplicationConfiguration = id.ApplicationConfiguration
type ApplicationProfileHistory = id.ApplicationProfileHistory
type ApplicationCatalog = id.ApplicationCatalog

type ImageType = id.ImageType
type VideoType = id.VideoType
type Stage = id.Stage

const (
	StageDevelopment       = id.StageDevelopment
	StageTesting           = id.StageTesting
	StagePilot             = id.StagePilot
	StageProduction        = id.StageProduction
	ObjectStatusAvailable  = id.ObjectStatusAvailable
	ObjectStatusProcessing = id.ObjectStatusProcessing
	ObjectStatusError      = id.ObjectStatusError
	ObjectStatusPending    = id.ObjectStatusPending
)

var EntidadesAutoMigrate = []interface{}{
	&StorageObject{},             // base para Image, Video
	&Image{},                     // depende de StorageObject
	&Video{},                     // depende de StorageObject
	&IntegrationType{},           // base para Configuracao, CatalogoApplication
	&TerminalModel{},             // base para Configuracao, CatalogoApplication
	&Application{},               // base para ApplicationProfileHistory
	&CategoryType{},              // base para Category
	&Category{},                  // depende de CategoryType
	&ApplicationContact{},        // depende de ApplicationProfileHistory
	&ApplicationDetail{},         // depende de ApplicationProfileHistory
	&ApplicationVersionHistory{}, // depende de Application
	&ApplicationConfiguration{},  // depende de IntegrationType, TerminalModel, Application
	&ApplicationProfileHistory{}, // depende de ApplicationContact, ApplicationDetail, Category, ApplicationConfiguration
	&ApplicationCatalog{},        // depende de ApplicationConfiguration + hook para FK
}

// PostMigratable interface para entidades que precisam de setup pós-migração
type PostMigratable interface {
	PostMigrate(db *gorm.DB) error
}

// RunPostMigrations executa PostMigrate para todas as entidades que implementam a interface
func runPostMigrations(db *gorm.DB) error {
	for _, entidade := range EntidadesAutoMigrate {
		if postMigratable, ok := entidade.(PostMigratable); ok {
			if err := postMigratable.PostMigrate(db); err != nil {
				return err
			}
		}
	}
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	for _, entidade := range EntidadesAutoMigrate {
		if err := db.AutoMigrate(entidade); err != nil {
			return err
		}
	}

	return runPostMigrations(db)
}
