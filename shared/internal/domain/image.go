package domain

type ImageType int16

const (
	ImageTypeIcon       ImageType = 0
	ImageTypeScreehshot ImageType = 1
	ImageTypeBanner     ImageType = 2
)

type Image struct {
	BaseModel
	StorageObjectId *int64        `gorm:"column:id_obj_armazenamento;not null;uniqueIndex" json:"storage_object_id"`
	StorageObject   StorageObject `gorm:"foreignKey:StorageObjectId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ImageType       ImageType     `gorm:"column:cod_tip_img;not null;check:cod_tip_img IN (0, 1, 2)" json:"image_type"`
}

func (Image) TableName() string { return "image" }
