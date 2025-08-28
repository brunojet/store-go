package domain

type Configuracao struct {
	BaseModel
	IdTipoIntegracao      int64                  `gorm:"column:id_tipo_integracao;not null;uniqueIndex:idx_cfg_unique,priority:1"`
	TipoIntegracao        TipoIntegracao         `gorm:"foreignKey:IdTipoIntegracao"`
	IdModeloTerminal      int64                  `gorm:"column:id_modelo_terminal;not null;uniqueIndex:idx_cfg_unique,priority:2"`
	ModeloTerminal        ModeloTerminal         `gorm:"foreignKey:IdModeloTerminal"`
	IdApp                 int64                  `gorm:"column:id_app;not null;uniqueIndex:idx_cfg_unique,priority:3"`
	Aplicativo            Aplicativo             `gorm:"foreignKey:IdApp"`
	VersoesAplicativo     []VersaoAplicativo     `gorm:"foreignKey:IdConfiguracao"`
	ConfiguracaoCadastros []ConfiguracaoCadastro `gorm:"foreignKey:IdConfiguracao"`
}

func (Configuracao) TableName() string { return "cfg" }
