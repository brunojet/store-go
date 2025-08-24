package domain

type Imagem struct {
	ID                int64              `gorm:"primaryKey;autoIncrement:false"` // PK sincronizado com Anexo
	Anexo             *Anexo             `gorm:"foreignKey:ID;references:ID;constraint:OnDelete:RESTRICT;belongsTo"`
	VersoesAplicativo []VersaoAplicativo `gorm:"foreignKey:IdImagem"`
}

func (Imagem) TableName() string { return "imagem" }
