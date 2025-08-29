package domain

type Aplicativo struct {
	BaseEntity
	Cadastros     []Cadastro     `gorm:"foreignKey:IdApp"`
	Configuracoes []Configuracao `gorm:"foreignKey:IdApp"`
}

func (Aplicativo) TableName() string { return "app" }
