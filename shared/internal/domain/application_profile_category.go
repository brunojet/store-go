package domain

import (
	"gorm.io/gorm"

	domain_tools "github.com/brunojet/store-go/shared/internal/domain_utils"
)

type ApplicationProfileHistoryCategory struct {
	ID         int64 `gorm:"primaryKey:priority:0;column:id"`
	CategoryID int64 `gorm:"primaryKey:priority:1;column:category_id"`

	ApplicationProfileHistory ApplicationProfileHistory `gorm:"-:migration;foreignKey:ID;references:ID"`
	Category                  Category                  `gorm:"-:migration;foreignKey:CategoryID;references:ID"`
}

func (ApplicationProfileHistoryCategory) TableName() string {
	return "application_profile_history_category"
}

// PostMigrate cria as constraints FK manualmente para a tabela de junção category
func (ApplicationProfileHistoryCategory) PostMigrate(db *gorm.DB) error {
	if err := domain_tools.CreateConstraints(
		db,
		"application_profile_history_category",
		"application_profile_history",
		"fk_app_profile_history_category_f0",
		[]string{"id"},
		"CASCADE",
		"RESTRICT",
	); err != nil {
		return err
	}

	if err := domain_tools.CreateConstraintsWithTargetFields(
		db,
		"application_profile_history_category",
		"category",
		"fk_app_profile_history_category_f1",
		[]string{"category_id"},
		[]string{"id"},
		"RESTRICT",
		"RESTRICT",
	); err != nil {
		return err
	}

	return nil
}
