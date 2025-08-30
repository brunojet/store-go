package domain

import "time"

type HistoricoPerfilAplicativo struct {
	BaseModel
	IdAplicativo            int64                    `gorm:"column:id_aplv"`
	Aplicativo              Aplicativo               `gorm:"foreignKey:IdAplicativo"`
	IdContatoAplicativo     int64                    `gorm:"column:id_cnt_aplv"`
	ContatoAplicativo       ContatoAplicativo        `gorm:"foreignKey:IdContatoAplicativo"`
	IdDetalheAplicativo     int64                    `gorm:"column:id_dtlh_aplv"`
	DetalheAplicativo       DetalheAplicativo        `gorm:"foreignKey:IdDetalheAplicativo"`
	ConfiguracoesAplicativo []ConfiguracaoAplicativo `gorm:"many2many:cfg_hist_pfl_aplv;joinForeignKey:id_hist_pfl_aplv;joinReferences:id_cfg_aplv"`
	Categorias              []Categoria              `gorm:"many2many:ctgr_hist_pfl_aplv;joinForeignKey:id_hist_pfl_aplv;joinReferences:id_ctgr"`
	CatalogoAplicativo      []CatalogoAplicativo     `gorm:"foreignKey:IdHistoricoPerfilAplicativo"`
	DthInicioRevisao        *time.Time               `gorm:"column:dth_ini_rvs"`
	DthInicioProducao       *time.Time               `gorm:"column:dth_ini_prd"`
	DthFim                  *time.Time               `gorm:"column:dth_fim"`
	MotivoFim               string                   `gorm:"column:motivo_fim;size:255"`
}

func (HistoricoPerfilAplicativo) TableName() string { return "hist_pfl_aplv" }
