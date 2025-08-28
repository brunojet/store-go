package domain

import (
	"gorm.io/gorm"
)

type Aplicativo struct {
	BaseEntity
	Cadastros     []Cadastro     `gorm:"foreignKey:IdApp"`
	Configuracoes []Configuracao `gorm:"foreignKey:IdApp"`
}

func (c *Aplicativo) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (Aplicativo) TableName() string { return "app" }
