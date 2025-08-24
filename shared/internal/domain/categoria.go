package domain

type Categoria struct {
	BaseEntity
	IdTipoCategoria int64          `gorm:"column:id_tipo_categoria;not null" json:"id_tipo_categoria"`
	TipoCategoria   *TipoCategoria `gorm:"foreignKey:IdTipoCategoria" json:"tipo_categoria,omitempty"`
	IdPai           *int64         `gorm:"column:id_pai" json:"id_pai,omitempty"`
	Pai             *Categoria     `gorm:"foreignKey:IdPai" json:"pai,omitempty"`
}

func (Categoria) TableName() string { return "cat" }
