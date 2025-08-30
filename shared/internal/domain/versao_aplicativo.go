package domain

import "time"

type VersaoAplicativo struct {
	BaseEntity
	IdConfiguracaoAplicativo int64                  `gorm:"column:id_cfg_aplv;not null;index"`
	ConfiguracaoAplicativo   ConfiguracaoAplicativo `gorm:"foreignKey:IdConfiguracaoAplicativo"`
	DthInicioPiloto          *time.Time             `gorm:"column:dth_ini_plt"`
	DthInicioProducao        *time.Time             `gorm:"column:dth_ini_prd"`
	DthFim                   *time.Time             `gorm:"column:dth_fim"`
	MotivoFim                string                 `gorm:"column:motivo_fim;size:255"`
	Tamanho                  int64                  `gorm:"column:tamanho;not null"`
	NomeVersao               string                 `gorm:"column:nome_versao;size:16;not null"`
	IdImagem                 int64                  `gorm:"column:id_img;index"`
	Imagem                   Imagem                 `gorm:"foreignKey:IdImagem;references:ID"`
}

func (VersaoAplicativo) TableName() string { return "vrs_app" }
