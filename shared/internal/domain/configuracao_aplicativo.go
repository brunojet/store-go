package domain

type ConfiguracaoAplicativo struct {
	BaseModel
	IdTipoIntegracao  int64              `gorm:"column:id_tip_itgr;not null;index;uniqueIndex:idx_cfg_unique,priority:1;index:idx_cfg_integracao_modelo,priority:1"`
	TipoIntegracao    TipoIntegracao     `gorm:"foreignKey:IdTipoIntegracao"`
	IdModeloTerminal  int64              `gorm:"column:id_mdl_trml;not null;index;uniqueIndex:idx_cfg_unique,priority:2;index:idx_cfg_integracao_modelo,priority:2"`
	ModeloTerminal    ModeloTerminal     `gorm:"foreignKey:IdModeloTerminal"`
	IdAplicativo      int64              `gorm:"column:id_aplv;not null;index;uniqueIndex:idx_cfg_unique,priority:3"`
	Aplicativo        Aplicativo         `gorm:"foreignKey:IdAplicativo"`
	VersoesAplicativo []VersaoAplicativo `gorm:"foreignKey:IdConfiguracaoAplicativo"`
}

func (ConfiguracaoAplicativo) TableName() string { return "cfg_aplv" }
