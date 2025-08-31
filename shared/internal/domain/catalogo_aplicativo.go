package domain

type CatalogoAplicativo struct {
	BaseModel
	IdConfiguracaoPerfilAplicativo int64                     `gorm:"not null;column:id_cfg_aplv;uniqueIndex:idx_ctlg_aplv_unique,priority:0"`
	ConfiguracaoPerfilAplicativo   ConfiguracaoAplicativo    `gorm:"foreignKey:IdConfiguracaoPerfilAplicativo"`
	IdEstagio                      int64                     `gorm:"not null;column:id_est;uniqueIndex:idx_ctlg_aplv_unique,priority:1"`
	Estagio                        Estagio                   `gorm:"foreignKey:IdEstagio"`
	IdHistoricoPerfilAplicativo    int64                     `gorm:"not null;column:id_hist_pfl_aplv;index"`
	HistoricoPerfilAplicativo      HistoricoPerfilAplicativo `gorm:"foreignKey:IdHistoricoPerfilAplicativo"`
	IdVersaoAplicativo             int64                     `gorm:"not null;column:id_vrs_aplv;index"`
	VersaoAplicativo               VersaoAplicativo          `gorm:"foreignKey:IdVersaoAplicativo"`
}

func (CatalogoAplicativo) TableName() string { return "ctlg_aplv" }
