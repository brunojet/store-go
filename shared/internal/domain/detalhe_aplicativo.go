package domain

type DetalheAplicativo struct {
	BaseModel
	Descricao string   `gorm:"column:descricao;size:255"`
	Imagens   []Imagem `gorm:"many2many:img_dtlh;joinForeignKey:id_dtlh_aplv;joinReferences:id_img"`
}

func (DetalheAplicativo) TableName() string { return "dtlh_aplv" }
