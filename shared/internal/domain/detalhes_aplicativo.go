package domain

import (
	"gorm.io/gorm"
)

type DetalhesAplicativo struct {
	BaseModel
	Descricao string          `gorm:"column:descricao;size:255"`
	Imagens   []ImagemDetalhe `gorm:"many2many:imagens_detalhes;"` // Relação N:N para reuso de imagens
}

type ImagemDetalhe struct {
	BaseModel
	IdImagem             int64              `gorm:"column:id_imagem;not null"`
	Imagem               Imagem             `gorm:"foreignKey:IdImagem"`
	IdDetalhesAplicativo int64              `gorm:"column:id_detalhes_aplicativo;not null"`
	DetalhesAplicativo   DetalhesAplicativo `gorm:"foreignKey:IdDetalhesAplicativo"`
	Descricao            string             `gorm:"column:descricao;size:255"`
}

func (DetalhesAplicativo) TableName() string { return "dtlh_app" }

func (d *DetalhesAplicativo) BeforeCreate(tx *gorm.DB) (err error) {
	// Adicione validações se necessário
	return nil
}
