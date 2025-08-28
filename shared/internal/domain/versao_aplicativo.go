package domain

type VersaoAplicativo struct {
	BaseEntity
	IdConfiguracao int64        `gorm:"column:id_configuracao;not null;index"`
	Configuracao   Configuracao `gorm:"foreignKey:IdConfiguracao"`
	Tamanho        int64        `gorm:"column:tamanho;not null"`
	NomeVersao     string       `gorm:"column:nome_versao;size:16;not null"`
	IdImagem       int64        `gorm:"column:id_imagem;index"`
	Imagem         Imagem       `gorm:"foreignKey:IdImagem;references:ID"`
}

func (VersaoAplicativo) TableName() string { return "versao_app" }
