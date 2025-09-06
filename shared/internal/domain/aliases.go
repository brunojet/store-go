package domain

import (
	infra "github.com/brunojet/infra-go/pkg/domain"
	"gorm.io/gorm"
)

type BaseModel = infra.BaseModel
type BaseEntity = infra.BaseEntity
type Anexo = infra.Anexo

type PostMigratable interface {
	PostMigrate(db *gorm.DB) error
}
