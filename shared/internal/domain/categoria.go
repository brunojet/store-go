package domain

type Categoria struct {
	BaseEntity
	IdTipoCategoria int64          `gorm:"column:id_tip_ctgr;not null" json:"id_tip_ctgr"`
	TipoCategoria   *TipoCategoria `gorm:"foreignKey:IdTipoCategoria" json:"tip_ctgr,omitempty"`
	IdPai           *int64         `gorm:"column:id_pai" json:"id_pai,omitempty"`
	Pai             *Categoria     `gorm:"foreignKey:IdPai" json:"pai,omitempty"`
}

func (Categoria) TableName() string { return "ctgr" }
