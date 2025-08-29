package domain

type CatalogoAplicativo struct {
	BaseModel
	IdConfiguracao     int64            `gorm:"not null;uniqueIndex:idx_catapp_unique,priority:0"`
	Configuracao       Configuracao     `gorm:"foreignKey:IdConfiguracao"`
	IdEstagio          int64            `gorm:"not null;uniqueIndex:idx_catapp_unique,priority:1"`
	Estagio            Estagio          `gorm:"foreignKey:IdEstagio"`
	IdCadastro         int64            `gorm:"not null;column:id_cadastro"`
	Cadastro           Cadastro         `gorm:"foreignKey:IdCadastro"`
	IdVersaoAplicativo int64            `gorm:"not null;index"`
	VersaoAplicativo   VersaoAplicativo `gorm:"foreignKey:IdVersaoAplicativo"`
	Ativo              bool
}

func (CatalogoAplicativo) TableName() string { return "cat_app" }
