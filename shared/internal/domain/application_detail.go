package domain

import (
	domain_tools "github.com/brunojet/store-go/shared/internal/domain_utils"
	"gorm.io/gorm"
)

type ApplicationDetail struct {
	BaseModel
	Description string  `gorm:"column:description;size:255"`
	Screenshots []Image `gorm:"many2many:application_detail_screenshots;joinForeignKey:ApplicationDetailID;joinReferences:ImageID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationDetail) TableName() string { return "application_detail" }

func (ApplicationDetail) PostMigrate(db *gorm.DB) error {
	if err := domain_tools.EnsureCascadeOnDelete(
		db,
		"application_detail_screenshots",
		"application_detail_id",
		"application_detail",
		"id",
	); err != nil {
		return err
	}

	return nil
}
