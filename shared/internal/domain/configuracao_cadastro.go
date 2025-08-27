package domain

type ConfiguracaoCadastro struct {
	BaseModel
	IdCadastro         int64        `gorm:"column:id_cadastro;not null"`
	DetalhesAplicativo Cadastro     `gorm:"foreignKey:IdCadastro"`
	IdConfiguracao     int64        `gorm:"column:id_configuracao;not null;index"`
	Configuracao       Configuracao `gorm:"foreignKey:IdConfiguracao"`
}

func (ConfiguracaoCadastro) TableName() string { return "cfg_cad" }
