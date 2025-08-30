package domain

type Aplicativo struct {
	BaseEntity
	Cadastros     []HistoricoPerfilAplicativo `gorm:"foreignKey:IdAplicativo"`
	Configuracoes []ConfiguracaoAplicativo    `gorm:"foreignKey:IdAplicativo"`
}

func (Aplicativo) TableName() string { return "aplv" }
