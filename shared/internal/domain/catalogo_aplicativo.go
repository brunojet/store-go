package domain

type CatalogoAplicativo struct {
	BaseModel
	IdTipoIntegracao   int64            `gorm:"not null;index;index:idx_catapp_unique,priority:1"`
	IdModeloTerminal   int64            `gorm:"not null;index;index:idx_catapp_unique,priority:2"`
	IdEstagio          int64            `gorm:"not null;index;index:idx_catapp_unique,priority:3"`
	IdApp              int64            `gorm:"not null;index;index:idx_catapp_unique,priority:4"`
	IdVersaoAplicativo int64            `gorm:"not null;index"`
	VersaoAplicativo   VersaoAplicativo `gorm:"foreignKey:IdVersaoAplicativo"`
	Estagio            EstagioCatalogo  `gorm:"foreignKey:IdEstagio"`
	Ativo              bool
}

func (CatalogoAplicativo) TableName() string { return "cat_app" }
