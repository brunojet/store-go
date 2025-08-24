package repo_test

import (
	"context"
	"fmt"
	"testing"

	repo "store-go/shared/internal/repo"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestModel é um tipo simples usado apenas nos testes genéricos do repositório.
type TestModel struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:128;not null"`
}

func (TestModel) TableName() string { return "test_models" }

func setupGenericDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}
	if err := db.AutoMigrate(&TestModel{}); err != nil {
		t.Fatalf("auto migrate test model failed: %v", err)
	}
	return db
}

func TestRepository_GenericCRUD(t *testing.T) {
	ctx := context.Background()
	db := setupGenericDB(t)
	r := repo.NewRepository[TestModel](db)

	m := &TestModel{Name: "alpha"}
	if err := r.Create(ctx, m); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if m.ID == 0 {
		t.Fatalf("expected ID set after create")
	}
	loaded, err := r.GetByID(ctx, m.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if loaded.Name != m.Name {
		t.Fatalf("expected name %s got %s", m.Name, loaded.Name)
	}
	// Update
	loaded.Name = "beta"
	if err := r.Update(ctx, loaded); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	reloaded, err := r.GetByID(ctx, m.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if reloaded.Name != "beta" {
		t.Fatalf("expected name beta got %s", reloaded.Name)
	}
	// Delete
	if err := r.Delete(ctx, reloaded); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	// reuse err variable from above - do not use := which would require a new var
	_, err = r.GetByID(ctx, m.ID)
	if err == nil {
		t.Fatalf("expected error getting deleted record")
	}
}

func TestRepository_ListWithParamsAndWithTx(t *testing.T) {
	ctx := context.Background()
	db := setupGenericDB(t)
	r := repo.NewRepository[TestModel](db)

	// create 30 items
	for i := 1; i <= 30; i++ {
		if err := r.Create(ctx, &TestModel{Name: fmt.Sprintf("item-%02d", i)}); err != nil {
			t.Fatalf("create item failed: %v", err)
		}
	}

	// Use ListWithParams with mod to filter items with odd numbers (simulate complex condition)
	params := &repo.ListParams{Order: "id asc", Offset: 5, Limit: 7}
	mod := func(q *gorm.DB) *gorm.DB {
		// use SQL expression to pick items where id % 2 = 1 (odd ids)
		return q.Where("(id % 2) = 1")
	}
	items, total, err := r.ListWithParams(ctx, params, mod)
	if err != nil {
		t.Fatalf("ListWithParams failed: %v", err)
	}
	if total == 0 {
		t.Fatalf("expected total > 0 but got 0")
	}
	if len(items) > params.Limit {
		t.Fatalf("expected at most %d items but got %d", params.Limit, len(items))
	}

	// Test WithTx rollback behavior: intentionally return error to force rollback
	origCount := int64(0)
	if err := db.Model(&TestModel{}).Count(&origCount).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	err = r.WithTx(ctx, func(txRepo *repo.Repository[TestModel]) error {
		if err := txRepo.Create(ctx, &TestModel{Name: "tx-item"}); err != nil {
			return err
		}
		// return error to rollback
		return fmt.Errorf("force rollback")
	})
	if err == nil {
		t.Fatalf("expected error from tx function")
	}
	var afterCount int64
	if err := db.Model(&TestModel{}).Count(&afterCount).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if afterCount != origCount {
		t.Fatalf("expected rollback so count unchanged (%d vs %d)", origCount, afterCount)
	}
}
