package domain

type Imagem struct {
	BaseModel
	IdAnexo           int64              `gorm:"column:id_anexo;not null;index"`
	Anexo             Anexo              `gorm:"foreignKey:IdAnexo;references:ID"`
	VersoesAplicativo []VersaoAplicativo `gorm:"foreignKey:IdImagem"`
}

func (Imagem) TableName() string { return "imagem" }
