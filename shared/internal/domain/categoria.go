package domain

type Categoria struct {
	BaseEntity
	IdTipoCategoria int64         `gorm:"column:id_tipo_categoria;not null"`
	TipoCategoria   TipoCategoria `gorm:"foreignKey:IdTipoCategoria"`
	IdPai           *int64        `gorm:"column:id_pai"`
	Pai             *Categoria    `gorm:"foreignKey:IdPai"`
}

func (Categoria) TableName() string { return "cat" }
