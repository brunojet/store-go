package domain

type Cadastro struct { //Deveria ser uma many to one?
	BaseModel
	IdAplicativo          *int64                 `gorm:"column:id_app"`
	Aplicativo            Aplicativo             `gorm:"foreignKey:IdAplicativo"`
	IdContato             *int64                 `gorm:"column:id_contato"`
	Contato               Contato                `gorm:"foreignKey:IdContato"`
	ConfiguracaoCadastros []ConfiguracaoCadastro `gorm:"foreignKey:IdCadastro"`
	CatalogoAplicativos   []CatalogoAplicativo   `gorm:"foreignKey:IdCadastro"`
}

func (Cadastro) TableName() string { return "cad" }
