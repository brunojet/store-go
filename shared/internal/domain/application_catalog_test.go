package domain_test

import (
	"database/sql"
	"testing"

	"github.com/brunojet/store-go/shared/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

// Remove invalid type aliases for constants; use domain.Stage and domain.StageDevelopment, etc. directly in tests.

var (
	ApplicationCatalogPkKeys = []string{"integration_type_id", "terminal_model_id", "stage", "application_id"}
)

const (
	ApplicationCatalogFkName         = "fk_application_catalog_configuration"
	ApplicationCatalogCreateTableSQL = `CREATE TABLE application_catalog (
	application_id INTEGER NOT NULL,
	ativo numeric,
	integration_type_id INTEGER NOT NULL,
	terminal_model_id INTEGER NOT NULL,
	stage INTEGER NOT NULL,
	PRIMARY KEY (integration_type_id, terminal_model_id, stage, application_id),
	CONSTRAINT fk_application_catalog_configuration
		FOREIGN KEY (application_id, integration_type_id, terminal_model_id)
		REFERENCES application_configuration(application_id, integration_type_id, terminal_model_id)
		ON DELETE RESTRICT ON UPDATE RESTRICT
	)`
	ApplicationCatalogAlterTableSQL = `ALTER TABLE application_catalog 
		   ADD CONSTRAINT fk_application_catalog_configuration 
		   FOREIGN KEY (application_id, integration_type_id, terminal_model_id) 
		   REFERENCES application_configuration(application_id, integration_type_id, terminal_model_id) 
		   ON DELETE RESTRICT ON UPDATE RESTRICT`
)

func TestApplicationCatalog_TableName(t *testing.T) {
	var v domain.ApplicationCatalog
	expected := "application_catalog"
	if got := v.TableName(); got != expected {
		t.Fatalf("TableName() = %v, want %v", got, expected)
	}
}

func setupTestDBWithCatalog(t *testing.T) *gorm.DB {
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open sqlite: %v", err)
	}

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&domain.ApplicationConfiguration{}, &domain.ApplicationCatalog{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	catalog := domain.ApplicationCatalog{}
	err = catalog.PostMigrate(db, ApplicationCatalogPkKeys, ApplicationCatalogFkName, ApplicationCatalogAlterTableSQL, ApplicationCatalogCreateTableSQL)
	if err != nil {
		t.Fatalf("PostMigrate failed: %v", err)
	}

	// Garante que o SQLite respeite as foreign keys
	db.Exec("PRAGMA foreign_keys = ON")

	return db
}

func TestApplicationCatalog_PostMigrate(t *testing.T) {
	db := setupTestDBWithCatalog(t)
	_ = db // apenas valida o setup
}

func TestApplicationCatalog_StageConstants(t *testing.T) {
	// Testar valores das constantes de estágio
	tests := []struct {
		stage    domain.Stage
		expected domain.Stage
		name     string
	}{
		{domain.StageDevelopment, 0, "StageDevelopment"},
		{domain.StageTesting, 10, "StageTesting"},
		{domain.StagePilot, 20, "StagePilot"},
		{domain.StageProduction, 30, "StageProduction"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stage != tt.expected {
				t.Fatalf("%s = %v, want %v", tt.name, tt.stage, tt.expected)
			}
		})
	}
}

// TestApplicationCatalog_ForeignKeyEnforcement tests FK constraints and deletion order.
func TestApplicationCatalog_ForeignKeyEnforcement(t *testing.T) {
	db := setupTestDBWithCatalog(t)

	// 1. Inserir dependências na ordem correta
	app := domain.Application{BaseEntity: domain.BaseEntity{Nome: "App Teste"}}
	it := domain.IntegrationType{BaseEntity: domain.BaseEntity{Nome: "Tipo Integracao Teste"}}
	tm := domain.TerminalModel{BaseEntity: domain.BaseEntity{Nome: "Terminal Teste"}}
	if err := db.Create(&app).Error; err != nil {
		t.Fatalf("Failed to insert Application: %v", err)
	}
	if err := db.Create(&it).Error; err != nil {
		t.Fatalf("Failed to insert IntegrationType: %v", err)
	}
	if err := db.Create(&tm).Error; err != nil {
		t.Fatalf("Failed to insert TerminalModel: %v", err)
	}
	acfg := domain.ApplicationConfiguration{
		ApplicationId:     app.ID,
		IntegrationTypeId: it.ID,
		TerminalModelId:   tm.ID,
	}
	if err := db.Create(&acfg).Error; err != nil {
		t.Fatalf("Failed to insert ApplicationConfiguration: %v", err)
	}
	acat := domain.ApplicationCatalog{
		IntegrationTypeId: it.ID,
		TerminalModelId:   tm.ID,
		Stage:             domain.StageDevelopment,
		ApplicationId:     app.ID,
	}
	if err := db.Create(&acat).Error; err != nil {
		t.Fatalf("Failed to insert ApplicationCatalog: %v", err)
	}

	// 2. Tentar excluir ApplicationConfiguration (deve falhar por FK)
	err := db.Delete(&domain.ApplicationConfiguration{}, "application_id = ? AND integration_type_id = ? AND terminal_model_id = ?", acfg.ApplicationId, acfg.IntegrationTypeId, acfg.TerminalModelId).Error
	if err == nil {
		t.Fatal("Expected FK violation when deleting ApplicationConfiguration, but got no error")
	}

	// 3. Tentar excluir Application (deve falhar por FK)
	err = db.Delete(&app).Error
	if err == nil {
		t.Fatal("Expected FK violation when deleting Application, but got no error")
	}

	// 4. Agora deletar na ordem correta
	// Primeiro remover domain.ApplicationCatalog
	if err := db.Delete(&domain.ApplicationCatalog{}, "integration_type_id = ? AND terminal_model_id = ? AND application_id = ?", acat.IntegrationTypeId, acat.TerminalModelId, acat.ApplicationId).Error; err != nil {
		t.Fatalf("Failed to delete domain.ApplicationCatalog: %v", err)
	}

	// Depois ApplicationConfiguration
	if err := db.Delete(&domain.ApplicationConfiguration{}, "application_id = ? AND integration_type_id = ? AND terminal_model_id = ?", acfg.ApplicationId, acfg.IntegrationTypeId, acfg.TerminalModelId).Error; err != nil {
		t.Fatalf("Failed to delete ApplicationConfiguration: %v", err)
	}
	// Depois TerminalModel
	if err := db.Delete(&domain.TerminalModel{}, "id = ?", tm.ID).Error; err != nil {
		t.Fatalf("Failed to delete TerminalModel: %v", err)
	}
	// Depois IntegrationType
	if err := db.Delete(&domain.IntegrationType{}, "id = ?", it.ID).Error; err != nil {
		t.Fatalf("Failed to delete IntegrationType: %v", err)
	}
	// Por fim Application
	if err := db.Delete(&domain.Application{}, "id = ?", app.ID).Error; err != nil {
		t.Fatalf("Failed to delete Application: %v", err)
	}
}
