package repo_test

import (
	"testing"

	"github.com/brunojet/store-go/infra/pkg/domain"
	pkgrepo "github.com/brunojet/store-go/infra/pkg/repo"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func withDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable fk: %v", err)
	}
	if err := db.AutoMigrate(&domain.Anexo{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestRepository_NewRepository(t *testing.T) {
	db := withDB(t)
	repo := pkgrepo.NewRepository[domain.Anexo](db)
	if repo.DB() != db {
		t.Errorf("expected db to be %v, got %v", db, repo.DB())
	}
}

func TestRepository_NewAnexoRepo(t *testing.T) {
	db := withDB(t)
	repo := pkgrepo.NewAnexoRepo(db)
	if repo.DB() != db {
		t.Errorf("expected db to be %v, got %v", db, repo.DB())
	}
}
