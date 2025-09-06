package domain

type ApplicationDetail struct {
	BaseModel
	Descricao string `gorm:"column:descricao;size:255"`
}

func (ApplicationDetail) TableName() string { return "application_detail" }
