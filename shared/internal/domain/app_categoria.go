package domain

type AppCategoria struct {
	BaseModel
	IdApp       int64     `gorm:"column:id_app;not null;index"`
	IdCategoria int64     `gorm:"column:id_categoria;not null;index"`
	Categoria   Categoria `gorm:"foreignKey:IdCategoria"`
}

func (AppCategoria) TableName() string { return "app_cat" }
