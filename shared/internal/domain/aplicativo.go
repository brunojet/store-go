package domain

type Aplicativo struct {
	BaseEntity
	HistoricoPerfilAplicativo []HistoricoPerfilAplicativo `gorm:"foreignKey:IdAplicativo"`
	ConfiguracoesAplicativo   []ConfiguracaoAplicativo    `gorm:"foreignKey:IdAplicativo"`
}

func (Aplicativo) TableName() string { return "aplv" }
