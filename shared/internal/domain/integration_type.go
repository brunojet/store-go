package domain

type IntegrationType struct {
	BaseEntity
}

func (IntegrationType) TableName() string { return "integration_type" }
