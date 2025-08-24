package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// --- shared test models and helpers ---
type TestModel struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:128;not null"`
}

func (TestModel) TableName() string { return "test_models" }

type TestModelPkg struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:128;not null"`
}

func (TestModelPkg) TableName() string { return "test_models_pkg" }

type ParentPkg struct {
	ID    int64 `gorm:"primaryKey;autoIncrement"`
	Name  string
	Child *ChildPkg `gorm:"foreignKey:ParentID"`
}

type ChildPkg struct {
	ID       int64
	ParentID int64
	Value    string
}

type MissingModel struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:128;not null"`
}

func (MissingModel) TableName() string { return "missing_models" }

func setupDBPkg(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}
	return db
}

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

// --- tests from previous internals file ---
func TestApplyPreloadsViaGetByID(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	if err := db.AutoMigrate(&ParentPkg{}, &ChildPkg{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	r := NewRepository[ParentPkg](db)

	p := &ParentPkg{Name: "parent-1"}
	if err := r.Create(ctx, p); err != nil {
		t.Fatalf("create parent failed: %v", err)
	}
	if err := db.Create(&ChildPkg{ParentID: p.ID, Value: "c1"}).Error; err != nil {
		t.Fatalf("create child failed: %v", err)
	}

	got, err := r.GetByID(ctx, p.ID, "Child")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Child == nil {
		t.Fatalf("expected Child preloaded")
	}
}

func TestListAndFindWithParams(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	if err := db.AutoMigrate(&TestModelPkg{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	r := NewRepository[TestModelPkg](db)

	for i := 1; i <= 15; i++ {
		if err := r.Create(ctx, &TestModelPkg{Name: fmt.Sprintf("it-%02d", i)}); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// List without mod should return all
	all, err := r.List(ctx, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(all) != 15 {
		t.Fatalf("expected 15 items got %d", len(all))
	}

	// FindWithParams should return total 15 and limited results
	params := &ListParams{Limit: 5, Offset: 2}
	items, total, err := r.FindWithParams(ctx, params, nil)
	if err != nil {
		t.Fatalf("FindWithParams failed: %v", err)
	}
	if total != 15 {
		t.Fatalf("expected total 15 got %d", total)
	}
	if len(items) != 5 {
		t.Fatalf("expected 5 items got %d", len(items))
	}
}

func TestDBMethod(t *testing.T) {
	db := setupDBPkg(t)
	r := NewRepository[TestModelPkg](db)
	if r.DB() == nil {
		t.Fatalf("DB() returned nil")
	}
}

func TestListWithModFilters(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	if err := db.AutoMigrate(&TestModelPkg{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	r := NewRepository[TestModelPkg](db)
	names := []string{"keep-1", "drop-1", "keep-2"}
	for _, n := range names {
		if err := r.Create(ctx, &TestModelPkg{Name: n}); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	// use mod to filter names starting with keep
	mod := func(q *gorm.DB) *gorm.DB {
		return q.Where("name LIKE ?", "keep-%")
	}
	items, err := r.List(ctx, mod)
	if err != nil {
		t.Fatalf("List with mod failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 filtered items got %d", len(items))
	}
}

func TestListWithParams_ModAndPreloads(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	if err := db.AutoMigrate(&ParentPkg{}, &ChildPkg{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	r := NewRepository[ParentPkg](db)

	p1 := &ParentPkg{Name: "p-a"}
	p2 := &ParentPkg{Name: "p-b"}
	if err := r.Create(ctx, p1); err != nil {
		t.Fatalf("create p1: %v", err)
	}
	if err := r.Create(ctx, p2); err != nil {
		t.Fatalf("create p2: %v", err)
	}
	if err := db.Create(&ChildPkg{ParentID: p1.ID, Value: "c-a"}).Error; err != nil {
		t.Fatalf("create child: %v", err)
	}

	params := &ListParams{Preloads: []string{"Child"}, Limit: 10}
	mod := func(q *gorm.DB) *gorm.DB { return q.Where("name = ?", "p-a") }
	items, total, err := r.ListWithParams(ctx, params, mod)
	if err != nil {
		t.Fatalf("ListWithParams failed: %v", err)
	}
	if total == 0 {
		t.Fatalf("expected total > 0")
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item got %d", len(items))
	}
	if items[0].Child == nil {
		t.Fatalf("expected child preloaded")
	}
}

func TestListReturnsEmptyWhenNoRows(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	if err := db.AutoMigrate(&TestModelPkg{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	r := NewRepository[TestModelPkg](db)

	// no inserts
	items, err := r.List(ctx, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items got %d", len(items))
	}
}

func TestListWithParamsReturnsEmptyWhenNoRows(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	if err := db.AutoMigrate(&TestModelPkg{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	r := NewRepository[TestModelPkg](db)

	params := &ListParams{Limit: 10, Offset: 0}
	items, total, err := r.ListWithParams(ctx, params, nil)
	if err != nil {
		t.Fatalf("ListWithParams failed: %v", err)
	}
	if total != 0 {
		t.Fatalf("expected total 0 got %d", total)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items got %d", len(items))
	}
}

func TestListReturnsErrorWhenTableMissing(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	// Intencionalmente NÃO chamamos AutoMigrate(&MissingModel{})
	r := NewRepository[MissingModel](db)

	_, err := r.List(ctx, nil)
	if err == nil {
		t.Fatalf("expected error when table is missing, got nil")
	}
}

func TestListWithParamsReturnsErrorWhenTableMissing(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	// Intencionalmente NÃO chamamos AutoMigrate(&MissingModel{})
	r := NewRepository[MissingModel](db)

	params := &ListParams{Limit: 10, Offset: 0}
	_, _, err := r.ListWithParams(ctx, params, nil)
	if err == nil {
		t.Fatalf("expected error from ListWithParams when table is missing, got nil")
	}
}

func TestListWithParamsCountFailsWithInvalidJoin(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	// Garantir tabela do modelo existe para que o erro venha do JOIN inexistente
	if err := db.AutoMigrate(&TestModelPkg{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	r := NewRepository[TestModelPkg](db)

	// inserir um registro para tornar a consulta válida se o JOIN existisse
	if err := r.Create(ctx, &TestModelPkg{Name: "x"}); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// mod aplica um JOIN para uma tabela que não existe, causando erro no COUNT
	mod := func(q *gorm.DB) *gorm.DB {
		return q.Joins("JOIN table_does_not_exist ON table_does_not_exist.id = test_models_pkg.id")
	}

	_, _, err := r.ListWithParams(ctx, &ListParams{Limit: 10}, mod)
	if err == nil {
		t.Fatalf("expected error from ListWithParams count when join references missing table, got nil")
	}
}

func TestListWithParamsFindFailsWhenPreloadMissingTable(t *testing.T) {
	ctx := context.Background()
	db := setupDBPkg(t)
	// Migrar apenas a tabela pai (ParentPkg), mas NÃO migrar ChildPkg
	if err := db.AutoMigrate(&ParentPkg{}); err != nil {
		t.Fatalf("migrate parent failed: %v", err)
	}
	r := NewRepository[ParentPkg](db)

	// inserir um pai válido
	p := &ParentPkg{Name: "parent-x"}
	if err := r.Create(ctx, p); err != nil {
		t.Fatalf("create parent failed: %v", err)
	}

	params := &ListParams{Preloads: []string{"Child"}, Limit: 10}
	_, _, err := r.ListWithParams(ctx, params, nil)
	if err == nil {
		t.Fatalf("expected error from ListWithParams Find when preload table missing, got nil")
	}
}

// --- tests from previous generic file ---
func TestRepository_GenericCRUD(t *testing.T) {
	ctx := context.Background()
	db := setupGenericDB(t)
	r := NewRepository[TestModel](db)

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
	r := NewRepository[TestModel](db)

	// create 30 items
	for i := 1; i <= 30; i++ {
		if err := r.Create(ctx, &TestModel{Name: fmt.Sprintf("item-%02d", i)}); err != nil {
			t.Fatalf("create item failed: %v", err)
		}
	}

	// Use ListWithParams with mod to filter items with odd numbers (simulate complex condition)
	params := &ListParams{Order: "id asc", Offset: 5, Limit: 7}
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
	err = r.WithTx(ctx, func(txRepo *Repository[TestModel]) error {
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

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}
	// migrate only test domain tables used here
	if err := db.AutoMigrate(&AnexoTest{}, &ImagemTest{}, &VersaoAplicativoTest{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	return db
}

// --- lightweight local test-domain types ---
type AnexoTest struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	Nome          string `gorm:"size:128;not null"`
	TipoMime      string `gorm:"size:64;not null"`
	MD5           string `gorm:"size:32;not null"`
	Tamanho       int64  `gorm:"not null"`
	Armazenamento string `gorm:"size:256;not null"`
	Caminho       string `gorm:"size:256;not null"`
	Presente      bool   `gorm:"default:false"`
}

func (AnexoTest) TableName() string { return "anexo" }

type ImagemTest struct {
	ID    int64     `gorm:"primaryKey;autoIncrement:false"`
	Anexo AnexoTest `gorm:"foreignKey:ID;references:ID;constraint:OnDelete:RESTRICT;belongsTo"`
}

func (ImagemTest) TableName() string { return "imagem" }

type VersaoAplicativoTest struct {
	ID       int64 `gorm:"primaryKey;autoIncrement"`
	IdImagem int64 `gorm:"column:id_imagem;index"`
}

func (VersaoAplicativoTest) TableName() string { return "versao_app" }

// test creating many versões referencing imagens (uses local test-domain types)

func TestListWithParams_PaginationAndMod(t *testing.T) {
	ctx := context.Background()
	db := setupDB(t)

	// create 5 anexo+imagem pairs because Imagem.ID is FK -> Anexo.ID
	for i := int64(1); i <= 5; i++ {
		an := &AnexoTest{
			ID:            i,
			Nome:          fmt.Sprintf("anexo-%d", i),
			TipoMime:      "image/png",
			MD5:           fmt.Sprintf("md5-%d", i),
			Tamanho:       100,
			Armazenamento: "local",
			Caminho:       fmt.Sprintf("a/%d", i),
			Presente:      true,
		}
		if err := db.Create(an).Error; err != nil {
			t.Fatalf("failed create anexo: %v", err)
		}
		im := &ImagemTest{ID: i}
		if err := db.Create(im).Error; err != nil {
			t.Fatalf("failed create imagem: %v", err)
		}
	}
	// create 25 versions: id_imagem cycles 1..5
	for i := int64(1); i <= 25; i++ {
		v := &VersaoAplicativoTest{
			IdImagem: ((i-1)%5 + 1),
		}
		if err := db.Create(v).Error; err != nil {
			t.Fatalf("failed create versao: %v", err)
		}
	}

	versRepo := NewRepository[VersaoAplicativoTest](db)
	// want to select versoes where id_imagem IN (2,3)
	params := &ListParams{Order: "id asc", Offset: 5, Limit: 7}
	mod := func(q *gorm.DB) *gorm.DB {
		return q.Where("id_imagem IN ?", []int64{2, 3})
	}
	items, total, err := versRepo.ListWithParams(ctx, params, mod)
	if err != nil {
		t.Fatalf("ListWithParams failed: %v", err)
	}
	if total == 0 {
		t.Fatalf("expected total > 0 for filtered versoes")
	}
	// ensure returned count respects limit (or less if fewer rows)
	if len(items) > params.Limit {
		t.Fatalf("expected at most %d items but got %d", params.Limit, len(items))
	}
}
