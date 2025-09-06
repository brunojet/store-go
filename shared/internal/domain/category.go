package domain

type Category struct {
	BaseEntity
	CategoryTypeId int64        `gorm:"column:category_type_id;not null;index" json:"category_type_id"`
	CategoryType   CategoryType `gorm:"foreignKey:CategoryTypeId" json:"category_type,omitempty"`
	ParentId       *int64       `gorm:"column:parent_id;index" json:"parent_id,omitempty"`
	Parent         *Category    `gorm:"foreignKey:ParentId" json:"parent,omitempty"`
}

func (Category) TableName() string { return "category" }
