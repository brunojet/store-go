package domain

type CategoryType struct {
	BaseEntity
}

func (CategoryType) TableName() string { return "category_type" }
