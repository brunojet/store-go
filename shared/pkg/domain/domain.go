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
type Anexo = infra.Anexo
type CategoryType = id.CategoryType
type Category = id.Category
type IntegrationType = id.IntegrationType
type TerminalModel = id.TerminalModel
type ApplicationContact = id.ApplicationContact
type Application = id.Application
type ApplicationDetail = id.ApplicationDetail
type ApplicationConfiguration = id.ApplicationConfiguration
type ApplicationProfileHistory = id.ApplicationProfileHistory
type ApplicationVersion = id.ApplicationVersionHistory
type Estagio = id.Estagio
type ApplicationCatalog = id.ApplicationCatalog
type StorageObject = id.StorageObject
type Image = id.Image
type ImageType = id.ImageType
type Video = id.Video
type VideoType = id.VideoType

var EntidadesAutoMigrate = []interface{}{
	&StorageObject{},            // base para Image, Video
	&Image{},                    // depende de StorageObject
	&Video{},                    // depende de StorageObject
	&IntegrationType{},          // base para Configuracao, CatalogoApplication
	&TerminalModel{},            // base para Configuracao, CatalogoApplication
	&Application{},              // base para AppCategory, Configuracao, Cadastro, CatalogoApplication
	&ApplicationConfiguration{}, // depende de IntegrationType, TerminalModel, Application
	&ApplicationCatalog{},       // depende de ApplicationConfiguration + hook para FK
	// &CategoryType{},              // base para Category
	// &Category{},                  // depende de CategoryType
	&ApplicationProfileHistory{}, // depende de ApplicationContact, ApplicationDetail, Category, ApplicationConfiguration
	&ApplicationVersion{},        // depende de Application
}

// PostMigratable interface para entidades que precisam de setup pós-migração
type PostMigratable interface {
	PostMigrate(db *gorm.DB) error
}

// RunPostMigrations executa PostMigrate para todas as entidades que implementam a interface
func RunPostMigrations(db *gorm.DB) error {
	for _, entidade := range EntidadesAutoMigrate {
		if postMigratable, ok := entidade.(PostMigratable); ok {
			if err := postMigratable.PostMigrate(db); err != nil {
				return err
			}
		}
	}
	return nil
}
