package domain

type Imagem struct {
	BaseModel
	IdAnexo           int64              `gorm:"column:id_anexo"`
	Anexo             Anexo              `gorm:"foreignKey:IdAnexo"`
	VersoesAplicativo []VersaoAplicativo `gorm:"foreignKey:IdImagem"`
}

func (Imagem) TableName() string { return "imagem" }
