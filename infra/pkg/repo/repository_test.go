package repo_test

import (
	"context"
	"testing"

	internalrepo "github.com/brunojet/store-go/infra/internal/repo"
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

func TestPublicRepository_CRUD(t *testing.T) {
	db := withDB(t)
	repo := pkgrepo.NewRepository[domain.Anexo](db)
	ctx := context.Background()

	an := &domain.Anexo{Nome: "t1", TipoMime: "image/png", MD5: "m", Tamanho: 10, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := repo.Create(ctx, an); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if an.ID == 0 {
		t.Fatalf("expected id set after create")
	}

	got, err := repo.GetByID(ctx, an.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Nome != an.Nome {
		t.Fatalf("unexpected name: %s", got.Nome)
	}

	// List should return the created entity
	list, err := repo.List(ctx, nil)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}

	// Delete and confirm
	if err := repo.Delete(ctx, an); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	l2, err := repo.List(ctx, nil)
	if err != nil {
		t.Fatalf("list after delete failed: %v", err)
	}
	if len(l2) != 0 {
		t.Fatalf("expected 0 items after delete, got %d", len(l2))
	}
}

func TestPublicRepository_Update(t *testing.T) {
	db := withDB(t)
	repo := pkgrepo.NewRepository[domain.Anexo](db)
	ctx := context.Background()

	an := &domain.Anexo{Nome: "orig", TipoMime: "image/png", MD5: "m", Tamanho: 10, Armazenamento: "loc", Caminho: "p", Presente: true}
	if err := repo.Create(ctx, an); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	an.Nome = "updated"
	if err := repo.Update(ctx, an); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	got, err := repo.GetByID(ctx, an.ID)
	if err != nil {
		t.Fatalf("get after update failed: %v", err)
	}
	if got.Nome != "updated" {
		t.Fatalf("expected name updated, got %s", got.Nome)
	}
}

func TestPublicRepository_ListWithParams_FindWithParams(t *testing.T) {
	db := withDB(t)
	repo := pkgrepo.NewRepository[domain.Anexo](db)
	ctx := context.Background()

	// create 3 items
	for i := 1; i <= 3; i++ {
		a := &domain.Anexo{Nome: "n", TipoMime: "image/png", MD5: "m", Tamanho: int64(i), Armazenamento: "loc", Caminho: "p", Presente: true}
		if err := repo.Create(ctx, a); err != nil {
			t.Fatalf("create %d failed: %v", i, err)
		}
	}

	// ListWithParams: default order and limit
	p := &internalrepo.ListParams{Limit: 2, Offset: 1}
	list, total, err := repo.ListWithParams(ctx, p, nil)
	if err != nil {
		t.Fatalf("ListWithParams failed: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 items in paged list, got %d", len(list))
	}

	// FindWithParams should behave the same
	f, ftot, err := repo.FindWithParams(ctx, p, nil)
	if err != nil {
		t.Fatalf("FindWithParams failed: %v", err)
	}
	if ftot != 3 {
		t.Fatalf("expected total 3 from FindWithParams, got %d", ftot)
	}
	if len(f) != 2 {
		t.Fatalf("expected 2 items from FindWithParams, got %d", len(f))
	}
}

func TestPublicRepository_WithTxAndDB(t *testing.T) {
	db := withDB(t)
	repo := pkgrepo.NewRepository[domain.Anexo](db)
	ctx := context.Background()

	// transaction rollback path: create then return error to rollback
	err := repo.WithTx(ctx, func(txRepo *pkgrepo.Repository[domain.Anexo]) error {
		a := &domain.Anexo{Nome: "tx", TipoMime: "image/png", MD5: "m", Tamanho: 1, Armazenamento: "loc", Caminho: "p", Presente: true}
		if err := txRepo.Create(ctx, a); err != nil {
			return err
		}
		// ensure visible inside tx via txRepo.DB()
		var cnt int64
		if err := txRepo.DB().Model(&domain.Anexo{}).Where("nome = ?", "tx").Count(&cnt).Error; err != nil {
			return err
		}
		if cnt != 1 {
			return nil // unexpected but continue to rollback
		}
		return gorm.ErrInvalidDB // force rollback
	})
	if err == nil {
		t.Fatalf("expected error from tx callback to cause rollback")
	}

	// after rollback the record should not exist
	var total int64
	if err := repo.DB().Model(&domain.Anexo{}).Where("nome = ?", "tx").Count(&total).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if total != 0 {
		t.Fatalf("expected 0 after rollback, got %d", total)
	}

	// transaction commit path
	err = repo.WithTx(ctx, func(txRepo *pkgrepo.Repository[domain.Anexo]) error {
		a := &domain.Anexo{Nome: "tx2", TipoMime: "image/png", MD5: "m", Tamanho: 2, Armazenamento: "loc", Caminho: "p", Presente: true}
		return txRepo.Create(ctx, a)
	})
	if err != nil {
		t.Fatalf("commit tx failed: %v", err)
	}
	// check via public DB()
	if err := repo.DB().Model(&domain.Anexo{}).Where("nome = ?", "tx2").Count(&total).Error; err != nil {
		t.Fatalf("count after commit failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 after commit, got %d", total)
	}
}
