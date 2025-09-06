package domain

type VideoType int16

const (
	VideoTypeShorts VideoType = 0
	VideoTypeReels  VideoType = 1
)

type Video struct {
	BaseModel
	StorageObjectId *int64        `gorm:"column:id_obj_armazenamento;not null;uniqueIndex" json:"storage_object_id"`
	StorageObject   StorageObject `gorm:"foreignKey:StorageObjectId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	VideoType       VideoType     `gorm:"column:cod_tip_vid;not null;check:cod_tip_vid IN (0, 1)" json:"video_type"`
}

func (Video) TableName() string { return "video" }
