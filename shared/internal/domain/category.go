package domain

type Category struct {
	BaseEntity
	CategoryTypeId int64  `gorm:"index;not null" json:"category_type_id"`
	ParentId       *int64 `gorm:"index" json:"parent_id,omitempty"`

	// Relations
	ApplicationProfiles *[]ApplicationProfileHistory `gorm:"many2many:application_profile_history_category;joinForeignKey:CategoryID;joinReferences:ProfileID" json:"application_profiles,omitempty"`
	CategoryType        *CategoryType                `gorm:"foreignKey:CategoryTypeId" json:"category_type,omitempty"`
	Parent              *Category                    `gorm:"foreignKey:ParentId" json:"parent,omitempty"`
}

// ;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT
func (Category) TableName() string { return "category" }
