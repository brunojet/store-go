package domain

type TipoCategoria struct {
	BaseEntity
}

func (TipoCategoria) TableName() string { return "tip_cat" }
