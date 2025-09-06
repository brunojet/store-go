package domain

import (
	"gorm.io/gorm"
)

// MigrateApplicationConfigWithCatalogFK ensures proper migration order and FK creation
// between ApplicationConfiguration (parent) and ApplicationCatalog (child).
func MigrateApplicationConfigWithCatalogFK(db *gorm.DB) error {
	// Step 1: Migrate parent table first (ApplicationConfiguration)
	if err := db.AutoMigrate(&ApplicationConfiguration{}); err != nil {
		return err
	}

	// Step 2: Migrate child table (ApplicationCatalog)
	if err := db.AutoMigrate(&ApplicationCatalog{}); err != nil {
		return err
	}

	// Step 3: Explicitly create the composite FK constraint
	// ApplicationCatalog -> ApplicationConfiguration
	if err := db.Migrator().CreateConstraint(&ApplicationCatalog{}, "ApplicationConfiguration"); err != nil {
		// Ignore if constraint already exists
		if !db.Migrator().HasConstraint(&ApplicationCatalog{}, "ApplicationConfiguration") {
			return err
		}
	}

	return nil
}
