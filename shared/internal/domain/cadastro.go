package domain

type Aplicativo struct {
	BaseEntity
	Contato
	AppCategorias []AppCategoria `gorm:"foreignKey:IdApp"`
	Configuraces  []Configuracao `gorm:"foreignKey:IdApp"`
}

type AppCategoria struct {
	BaseModel
	IdApp       int64     `gorm:"column:id_app;not null;index"`
	IdCategoria int64     `gorm:"column:id_categoria;not null;index"`
	Categoria   Categoria `gorm:"foreignKey:IdCategoria"`
}

type Configuracao struct {
	BaseModel
	IdTipoIntegracao      int64                  `gorm:"column:id_tipo_integracao;not null;uniqueIndex:idx_cfg_unique,priority:1"` // FK para TipoIntegracao
	IdModeloTerminal      int64                  `gorm:"column:id_modelo_terminal;not null;uniqueIndex:idx_cfg_unique,priority:2"` // FK para ModeloTerminal
	IdApp                 int64                  `gorm:"column:id_app;not null;uniqueIndex:idx_cfg_unique,priority:3"`             // FK para Aplicativo
	Aplicativo            Aplicativo             `gorm:"foreignKey:IdApp"`
	TipoIntegracao        TipoIntegracao         `gorm:"foreignKey:IdTipoIntegracao"`
	ModeloTerminal        ModeloTerminal         `gorm:"foreignKey:IdModeloTerminal"`
	VersoesAplicativo     []VersaoAplicativo     `gorm:"foreignKey:IdConfiguracao"`
	ConfiguracaoCadastros []ConfiguracaoCadastro `gorm:"foreignKey:IdConfiguracao"`
}

type Cadastro struct {
	BaseModel
	ConfiguracaoCadastros []ConfiguracaoCadastro `gorm:"foreignKey:IdCadastro"`
	VersoesAplicativo     []VersaoAplicativo     `gorm:"foreignKey:IdCadastro"`
}

type ConfiguracaoCadastro struct {
	BaseModel
	IdCadastro     int64        `gorm:"column:id_cadastro;not null"`
	Cadastro       Cadastro     `gorm:"foreignKey:IdCadastro"`
	IdConfiguracao int64        `gorm:"column:id_configuracao;not null;index"`
	Configuracao   Configuracao `gorm:"foreignKey:IdConfiguracao"`
}

type VersaoAplicativo struct {
	BaseEntity
	IdCadastro     int64        `gorm:"column:id_cadastro;not null"`
	Cadastro       Cadastro     `gorm:"foreignKey:IdCadastro"`
	IdConfiguracao int64        `gorm:"column:id_configuracao;not null;index"`
	Configuracao   Configuracao `gorm:"foreignKey:IdConfiguracao"`
	Tamanho        int64        `gorm:"column:tamanho;not null"`
	NomeVersao     string       `gorm:"column:nome_versao;size:16;not null"`
}

// Tabela de estágios de publicação/configuração
type EstagioCatalogo struct {
	BaseEntity
	Nome                string `gorm:"size:50;unique;not null"`
	Descricao           string `gorm:"size:255"`
	Ativo               bool
	CatalogosAplicativo []CatalogoAplicativo `gorm:"foreignKey:IdEstagio"`
}

// Tabela catálogo aplicativo
type CatalogoAplicativo struct {
	BaseModel
	IdApp              int64            `gorm:"not null;index"`                                    // FK para Aplicativo
	IdTipoIntegracao   int64            `gorm:"not null;index;index:idx_catapp_unique,priority:1"` // FK para TipoIntegracao
	IdModeloTerminal   int64            `gorm:"not null;index;index:idx_catapp_unique,priority:2"` // FK para ModeloTerminal
	IdEstagio          int64            `gorm:"not null;index;index:idx_catapp_unique,priority:3"` // FK para EstagioCatalogo
	IdVersaoAplicativo int64            `gorm:"not null;index"`                                    // FK para VersaoAplicativo
	VersaoAplicativo   VersaoAplicativo `gorm:"foreignKey:IdVersaoAplicativo"`
	Estagio            EstagioCatalogo  `gorm:"foreignKey:IdEstagio"`
	Ativo              bool
}

func (Aplicativo) TableName() string           { return "app" }
func (AppCategoria) TableName() string         { return "app_cat" }
func (Configuracao) TableName() string         { return "cfg" }
func (Cadastro) TableName() string             { return "cad" }
func (ConfiguracaoCadastro) TableName() string { return "cfg_cad" }
func (VersaoAplicativo) TableName() string     { return "versao_app" }
func (EstagioCatalogo) TableName() string      { return "est_cat" }
func (CatalogoAplicativo) TableName() string   { return "cat_app" }
