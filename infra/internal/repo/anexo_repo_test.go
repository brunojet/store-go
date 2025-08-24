package repo

import (
	"context"
	"testing"

	"github.com/brunojet/store-go/infra/pkg/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// imagemModel is a minimal local representation of the application Imagem
// table so tests can AutoMigrate and verify rows without importing the
// application domain Imagem type (keeps repo package decoupled).
type imagemModel struct {
	ID int64 `gorm:"primaryKey;autoIncrement:false"`
}

func (imagemModel) TableName() string { return "imagem" }

func withTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed open db: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed enable fk: %v", err)
	}
	// imagem is an application-level entity; use local imagemModel for migration
	if err := db.AutoMigrate(&domain.Anexo{}, &imagemModel{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	return db
}

func TestAnexoRepo_FindByNome(t *testing.T) {
	db := withTestDB(t)
	r := NewAnexoRepo(db)
	ctx := context.Background()

	a := &domain.Anexo{Nome: "n1", TipoMime: "image/png", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := r.Create(ctx, a); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	res, err := r.FindByNome(ctx, "n1")
	if err != nil {
		t.Fatalf("FindByNome failed: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
}

func TestAnexoRepo_CreateWith_SuccessAndChildCreated(t *testing.T) {
	db := withTestDB(t)
	r := NewAnexoRepo(db)
	ctx := context.Background()

	a := &domain.Anexo{Nome: "a1", TipoMime: "t", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := r.CreateWith(ctx, a, func(tx *gorm.DB, anexoID int64) error {
		return tx.Create(&imagemModel{ID: anexoID}).Error
	}); err != nil {
		t.Fatalf("CreateWith failed: %v", err)
	}

	// ensure child exists
	var c imagemModel
	if err := db.First(&c, a.ID).Error; err != nil {
		t.Fatalf("expected child row, got err: %v", err)
	}
}

func TestAnexoRepo_CreateWith_CallbackErrorCausesRollback(t *testing.T) {
	db := withTestDB(t)
	r := NewAnexoRepo(db)
	ctx := context.Background()

	a := &domain.Anexo{Nome: "a2", TipoMime: "t", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
	// callback returns error -> whole transaction should rollback (no anexo created)
	if err := r.CreateWith(ctx, a, func(tx *gorm.DB, anexoID int64) error {
		return tx.Exec("INVALID SQL").Error // force an error
	}); err == nil {
		t.Fatalf("expected error from CreateWith when callback fails")
	}

	// ensure anexo not present
	var got domain.Anexo
	if err := db.First(&got, "nome = ?", "a2").Error; err == nil {
		t.Fatalf("expected no anexo due to rollback, found one")
	}
}

func TestAnexoRepo_GetWith_ReturnsAnexoAndChildExists(t *testing.T) {
	db := withTestDB(t)
	// create parent + child directly via db to simulate existing rows
	an := &domain.Anexo{BaseModel: domain.BaseModel{ID: 1}, Nome: "withchild", TipoMime: "t", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := db.Create(an).Error; err != nil {
		t.Fatalf("create anexo: %v", err)
	}
	if err := db.Create(&imagemModel{ID: an.ID}).Error; err != nil {
		t.Fatalf("create child: %v", err)
	}

	r := NewAnexoRepo(db)
	got, err := r.GetWith(context.Background(), an.ID)
	if err != nil {
		t.Fatalf("GetWith failed: %v", err)
	}
	// verify returned pointer and basic fields
	if got == nil {
		t.Fatalf("expected non-nil result")
	}
	if got.Nome != "withchild" {
		t.Fatalf("expected Nome=withchild got=%s", got.Nome)
	}
	// ensure child exists in imagem table
	var img imagemModel
	if err := db.First(&img, an.ID).Error; err != nil {
		t.Fatalf("expected child row in imagem table: %v", err)
	}
}

func TestAnexoRepo_CreateWith_NoChildCallbackCreatesAnexo(t *testing.T) {
	db := withTestDB(t)
	r := NewAnexoRepo(db)
	ctx := context.Background()

	a := &domain.Anexo{Nome: "nochild", TipoMime: "t", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := r.CreateWith(ctx, a, nil); err != nil {
		t.Fatalf("CreateWith(nil) failed: %v", err)
	}

	// ensure anexo saved
	var got domain.Anexo
	if err := db.First(&got, "nome = ?", "nochild").Error; err != nil {
		t.Fatalf("expected anexo row created: %v", err)
	}
}

func TestAnexoRepo_CreateWith_CreateFailsReturnsError(t *testing.T) {
	// open a DB but do NOT AutoMigrate the anexo table so Create fails (table not found)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("enable fk: %v", err)
	}

	r := NewAnexoRepo(db)
	ctx := context.Background()
	a := &domain.Anexo{Nome: "willfail", TipoMime: "t", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}

	if err := r.CreateWith(ctx, a, nil); err == nil {
		t.Fatalf("expected CreateWith to fail due to missing table, but it succeeded")
	}
}
