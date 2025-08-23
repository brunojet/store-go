package domain

type ModeloTerminal struct {
	BaseEntity
}

type TipoIntegracao struct {
	BaseEntity
}

type TipoCategoria struct {
	BaseEntity
}

type Categoria struct {
	BaseEntity
	IdTipoCategoria int64         `gorm:"column:id_tipo_categoria;not null"` // Chave estrangeira para TipoCategoria (obrigatório)
	TipoCategoria   TipoCategoria `gorm:"foreignKey:IdTipoCategoria"`        // Relacionamento
	IdPai           *int64        `gorm:"column:id_pai"`                     // Chave estrangeira para categoria pai (opcional)
	Pai             *Categoria    `gorm:"foreignKey:IdPai"`                  // Auto-relacionamento
}

func (ModeloTerminal) TableName() string { return "mdl_trml" }
func (TipoIntegracao) TableName() string { return "tip_int" }
func (TipoCategoria) TableName() string  { return "tip_cat" }
func (Categoria) TableName() string      { return "cat" }
