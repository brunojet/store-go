package domain

import (
	"gorm.io/gorm"
)

type DetalhesAplicativo struct {
	BaseModel
	Descricao string   `gorm:"column:descricao;size:255"`
	Imagens   []Imagem `gorm:"many2many:imagem_detalhe;"` // Relação N:N para reuso de imagens
}

func (DetalhesAplicativo) TableName() string { return "dtlh_app" }

func (d *DetalhesAplicativo) BeforeCreate(tx *gorm.DB) (err error) {
	// Adicione validações se necessário
	return nil
}
