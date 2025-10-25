package domain_test

import (
	// "database/sql"

	"fmt"
	"testing"

	"github.com/brunojet/store-go/shared/internal/domain"
	domain_tools "github.com/brunojet/store-go/shared/internal/domain_utils"
	pdomain "github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"

	// _ "github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

func TestApplicationCatalog_TableName(t *testing.T) {
	var v domain.ApplicationCatalog
	expected := "application_catalog"
	if got := v.TableName(); got != expected {
		t.Fatalf("TableName() = %v, want %v", got, expected)
	}
}

func setupTestDBWithCatalog(t *testing.T) *gorm.DB {

	// --- AJUSTE TEMPORÁRIO PARA USAR POSTGRES REAL ---
	dsn := "root:p0o9i8u7@tcp(127.0.0.1:3306)/store?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatalf("Failed to connect to mysql: %v", err)
	}
	// --- FIM DO AJUSTE TEMPORÁRIO ---

	// --- CÓDIGO ORIGINAL SQLITE (comente para reverter) ---
	// sqlDB, err := sql.Open("sqlite", ":memory:")
	// if err != nil {
	//     t.Fatalf("Failed to open sqlite: %v", err)
	// }
	// db, err := gorm.Open(sqlite.New(sqlite.Config{
	//     Conn: sqlDB,
	// }), &gorm.Config{
	//     Logger: logger.Default.LogMode(logger.Info),
	// })
	// if err != nil {
	//     t.Fatalf("Failed to connect to database: %v", err)
	// }

	// Use the public AutoMigrate re-export which composes internal entities
	err = pdomain.AutoMigrate(db)

	if err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

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
		{domain.StageDevelopment, 10, "StageDevelopment"},
		{domain.StageTesting, 20, "StageTesting"},
		{domain.StagePilot, 30, "StagePilot"},
		{domain.StageProduction, 40, "StageProduction"},
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
	app := domain.Application{BaseEntity: domain.BaseEntity{Name: "App Teste"}}
	it := domain.IntegrationType{BaseEntity: domain.BaseEntity{Name: "Tipo Integracao Teste"}}
	tm := domain.TerminalModel{BaseEntity: domain.BaseEntity{Name: "Terminal Teste"}}
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

	profile := domain.ApplicationProfileHistory{
		ApplicationContact: &domain.ApplicationContact{
			BaseEntity: domain.BaseEntity{Name: "Contato Teste"},
			Site:       "https://site.com",
			Email:      "contato@site.com",
			Phone:      "11999999999",
		},
		ApplicationDetail: &domain.ApplicationDetail{
			Description: "Detalhe do perfil",
		},
		Categories: []domain.Category{
			{
				BaseEntity: domain.BaseEntity{Name: "Categoria 1"},
				CategoryType: &domain.CategoryType{
					BaseEntity: domain.BaseEntity{Name: "Tipo Categoria 1"},
				},
			},
			{
				BaseEntity: domain.BaseEntity{Name: "Categoria 2"},
				CategoryType: &domain.CategoryType{
					BaseEntity: domain.BaseEntity{Name: "Tipo Categoria 2"},
				},
			},
		},
	}
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("Failed to insert ApplicationProfileHistory: %v", err)
	}
	// Associa via método em lote (mesmo com 1 elemento)
	if err := db.Model(&profile).Association("ApplicationConfigurations").Append(&acfg); err != nil {
		t.Fatalf("Failed to associate ApplicationConfiguration(s) with profile: %v", err)
	}

	version := domain.ApplicationVersionHistory{
		BaseEntity:        domain.BaseEntity{Name: "Versao Teste"},
		ApplicationId:     app.ID,
		IntegrationTypeId: it.ID,
		TerminalModelId:   tm.ID,
		Size:              123456,
		VersionName:       "1.0.0",
		Image: &domain.Image{
			StorageObject: domain.StorageObject{
				Path:     uuid.NewString(),
				Name:     "objeto-teste",
				MimeType: "etag-teste",
			},
		},
	}

	if err := db.Create(&version).Error; err != nil {
		t.Fatalf("Failed to insert ApplicationVersionHistory: %v", err)
	}

	acat := domain.ApplicationCatalog{
		IntegrationTypeId:    it.ID,
		TerminalModelId:      tm.ID,
		Stage:                domain.StageDevelopment,
		ApplicationId:        app.ID,
		ApplicationProfileId: &profile.ID,
		ApplicationVersionId: &version.ID,
	}
	if err := db.Create(&acat).Error; err != nil {
		t.Fatalf("Failed to insert ApplicationCatalog: %v", err)
	}

	// 2. Tentar excluir ApplicationConfiguration (deve falhar por FK)
	err := db.Delete(&domain.ApplicationConfiguration{}, "application_id = ? AND integration_type_id = ? AND terminal_model_id = ?", acfg.ApplicationId, acfg.IntegrationTypeId, acfg.TerminalModelId).Error
	if err == nil {
		t.Fatal("Expected FK violation when deleting ApplicationConfiguration, but got no error")
	}

	var loadedApp domain.Application
	err = db.Preload("ApplicationConfigurations").
		Preload("ApplicationConfigurations.IntegrationType").
		Preload("ApplicationConfigurations.TerminalModel").
		Preload("ApplicationConfigurations.ApplicationProfiles").
		Preload("ApplicationConfigurations.ApplicationProfiles.Categories").
		Preload("ApplicationConfigurations.ApplicationProfiles.Categories.CategoryType").
		Preload("ApplicationConfigurations.ApplicationVersions").
		Preload("ApplicationConfigurations.ApplicationVersions.Image").
		First(&loadedApp, app.ID).Error
	if err != nil {
		t.Fatalf("Failed to recursively load Application with relations: %v", err)
	}

	b, err := domain_tools.MarshalIndentWithoutNulls(loadedApp, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal loadedApp: %v", err)
	}
	fmt.Println(string(b))

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

	if err := db.Delete(&version).Error; err != nil {
		t.Fatalf("Failed to delete ApplicationVersionHistory: %v", err)
	}

	if err := db.Delete(&version.Image).Error; err != nil {
		t.Fatalf("Failed to delete ApplicationVersionHistory: %v", err)
	}

	if err := db.Delete(&profile).Error; err != nil {
		t.Fatalf("Failed to delete ApplicationProfileHistory: %v", err)
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
