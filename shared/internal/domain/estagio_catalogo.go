package domain

type Estagio struct {
	BaseEntity
	CatalogosApplication []ApplicationCatalog `gorm:"foreignKey:IdEstagio"`
}

func (Estagio) TableName() string { return "est" }
