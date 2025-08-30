package domain

type Cadastro struct { //Deveria ser uma many to one?
	BaseModel
	IdApp                int64                `gorm:"column:id_app"`
	Aplicativo           Aplicativo           `gorm:"foreignKey:IdApp"`
	IdContato            int64                `gorm:"column:id_contato"`
	Contato              Contato              `gorm:"foreignKey:IdContato"`
	IdDetalhesAplicativo int64                `gorm:"column:id_detalhes_aplicativo"`
	DetalhesAplicativo   DetalhesAplicativo   `gorm:"foreignKey:IdDetalhesAplicativo"`
	Configuracoes        []Configuracao       `gorm:"many2many:cfg_cad;joinForeignKey:cadastro_id;joinReferences:configuracao_id"`
	Categorias           []Categoria          `gorm:"many2many:cat_cad;joinForeignKey:cadastro_id;joinReferences:categoria_id"`
	CatalogoAplicativos  []CatalogoAplicativo `gorm:"foreignKey:IdCadastro"`
}

func (Cadastro) TableName() string { return "cad" }
