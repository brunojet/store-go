package domain

type Application struct {
	BaseEntity
	ConfiguracoesApplication []ApplicationConfiguration `gorm:"foreignKey:ApplicationId"`
}

func (Application) TableName() string { return "application" }
