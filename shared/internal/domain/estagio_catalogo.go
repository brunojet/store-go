package domain

type Estagio struct {
	BaseEntity
	CatalogosAplicativo []CatalogoAplicativo `gorm:"foreignKey:IdEstagio"`
}

func (Estagio) TableName() string { return "est_cat" }
