package domain

type Cadastro struct {
	BaseModel
	ConfiguracaoCadastros []ConfiguracaoCadastro `gorm:"foreignKey:IdCadastro"`
	VersoesAplicativo     []VersaoAplicativo     `gorm:"foreignKey:IdCadastro"`
}

func (Cadastro) TableName() string { return "cad" }
