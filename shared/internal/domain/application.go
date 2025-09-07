package domain

type Application struct {
	BaseEntity
	ApplicationConfigurations []ApplicationConfiguration `gorm:"foreignKey:ApplicationId"`
}

func (Application) TableName() string { return "application" }
