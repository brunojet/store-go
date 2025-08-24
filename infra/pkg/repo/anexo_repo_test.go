package repo_test

import (
	"context"
	"testing"

	"github.com/brunojet/store-go/infra/pkg/domain"
	pkgrepo "github.com/brunojet/store-go/infra/pkg/repo"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func withPkgDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("enable fk: %v", err)
	}
	if err := db.AutoMigrate(&domain.Anexo{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestPkgAnexoRepo_CreateWith_NoChildAndFindByNome(t *testing.T) {
	db := withPkgDB(t)
	r := pkgrepo.NewAnexoRepo(db)
	ctx := context.Background()

	a := &domain.Anexo{Nome: "pkg1", TipoMime: "t", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := r.CreateWith(ctx, a, nil); err != nil {
		t.Fatalf("CreateWith(nil) failed: %v", err)
	}

	res, err := r.FindByNome(ctx, "pkg1")
	if err != nil {
		t.Fatalf("FindByNome failed: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
}

func TestPkgAnexoRepo_CreateWith_CallbackErrorRollsBack(t *testing.T) {
	db := withPkgDB(t)
	r := pkgrepo.NewAnexoRepo(db)
	ctx := context.Background()

	a := &domain.Anexo{Nome: "pkg2", TipoMime: "t", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := r.CreateWith(ctx, a, func(tx *gorm.DB, anexoID int64) error {
		return tx.Exec("INVALID SQL").Error
	}); err == nil {
		t.Fatalf("expected CreateWith callback error")
	}
	// ensure no row
	var cnt int64
	if err := db.Model(&domain.Anexo{}).Where("nome = ?", "pkg2").Count(&cnt).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if cnt != 0 {
		t.Fatalf("expected 0 rows after rollback, got %d", cnt)
	}
}

func TestPkgAnexoRepo_GetWith_Basic(t *testing.T) {
	db := withPkgDB(t)
	r := pkgrepo.NewAnexoRepo(db)
	ctx := context.Background()

	a := &domain.Anexo{Nome: "pkg3", TipoMime: "t", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := r.CreateWith(ctx, a, nil); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := r.GetWith(ctx, a.ID)
	if err != nil {
		t.Fatalf("GetWith failed: %v", err)
	}
	if got == nil || got.Nome != "pkg3" {
		t.Fatalf("unexpected get result: %+v", got)
	}
}
