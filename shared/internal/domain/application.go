package domain

type Application struct {
	BaseEntity
	ApplicationConfigurations []ApplicationConfiguration `gorm:"foreignKey:ApplicationId" json:"application_configurations"`
}

func (Application) TableName() string { return "application" }
