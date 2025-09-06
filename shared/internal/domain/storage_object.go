package domain

type ObjectStatus int16

const (
	ObjectStatusPending    ObjectStatus = 0
	ObjectStatusProcessing ObjectStatus = 1
	ObjectStatusAvailable  ObjectStatus = 2
	ObjectStatusError      ObjectStatus = 3
)

type StorageObject struct {
	BaseModel
	Path     string       `gorm:"column:path;not null;uniqueIndex,type:char(40)" json:"path"`
	Name     string       `gorm:"column:name;not null;type:varchar(40)" json:"name"`
	MimeType string       `gorm:"column:mime_type;not null;type:varchar(100)" json:"mime_type"`
	Status   ObjectStatus `gorm:"column:status;not null" json:"status"`
}

func (StorageObject) TableName() string { return "storage_object" }
