package domain

type ConfiguracaoAplicativo struct {
	BaseModel
	IdTipoIntegracao  int64              `gorm:"column:id_tip_itgr;not null;index;uniqueIndex:idx_cfg_unique,priority:0;index:idx_cfg_itg_mdl,priority:1"`
	TipoIntegracao    TipoIntegracao     `gorm:"foreignKey:IdTipoIntegracao"`
	IdTerminalModel   int64              `gorm:"column:id_mdl_trml;not null;index;uniqueIndex:idx_cfg_unique,priority:1;index:idx_cfg_itg_mdl,priority:0"`
	TerminalModel     TerminalModel      `gorm:"foreignKey:IdTerminalModel"`
	IdAplicativo      int64              `gorm:"column:id_aplv;not null;index;uniqueIndex:idx_cfg_unique,priority:2"`
	Aplicativo        Aplicativo         `gorm:"foreignKey:IdAplicativo"`
	VersoesAplicativo []VersaoAplicativo `gorm:"foreignKey:IdConfiguracaoAplicativo"`
}

func (ConfiguracaoAplicativo) TableName() string { return "cfg_aplv" }
