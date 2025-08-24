package domain

type EstagioCatalogo struct {
	BaseEntity
	CatalogosAplicativo []CatalogoAplicativo `gorm:"foreignKey:IdEstagio"`
}

func (EstagioCatalogo) TableName() string { return "est_cat" }
