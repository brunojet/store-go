package domain

type Cadastro struct { //Deveria ser uma many to one?
	BaseModel
	IdApp                 int64                  `gorm:"column:id_app"`
	Aplicativo            Aplicativo             `gorm:"foreignKey:IdApp"`
	IdContato             int64                  `gorm:"column:id_contato"`
	Contato               Contato                `gorm:"foreignKey:IdContato"`
	IdDetalhesAplicativo  int64                  `gorm:"column:id_detalhes_aplicativo"`
	DetalhesAplicativo    DetalhesAplicativo     `gorm:"foreignKey:IdDetalhesAplicativo"`
	ConfiguracaoCadastros []ConfiguracaoCadastro `gorm:"foreignKey:IdCadastro"`
	CatalogoAplicativos   []CatalogoAplicativo   `gorm:"foreignKey:IdCadastro"`
}

func (Cadastro) TableName() string { return "cad" }
