package domain

type CadastroCategoria struct {
	BaseModel
	IdCategoria int64     `gorm:"column:id_categoria;not null;index"`
	Categoria   Categoria `gorm:"foreignKey:IdCategoria"`
	IdCadastro  int64     `gorm:"column:id_cadastro;not null;index"`
	Cadastro    Cadastro  `gorm:"foreignKey:IdCadastro"`
}

func (CadastroCategoria) TableName() string { return "cad_cat" }
