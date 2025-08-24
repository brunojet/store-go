package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// local TestModel for package tests
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

// MissingModel existe em Go, mas NÃO é migrado para o DB neste teste.
// Isso deve provocar um erro no Find (ex: "no such table") quando o GORM
// tentar executar a query contra a tabela inexistente.
type MissingModel struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:128;not null"`
}

func (MissingModel) TableName() string { return "missing_models" }

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

// Força erro na contagem usando um JOIN para tabela inexistente.
// Como Count é executado antes do Find em ListWithParams, este teste
// deve provocar erro no Count e cobrir o branch de erro nessa linha.
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

// Provoca erro durante o q.Find em ListWithParams ao solicitar Preload
// de relação cuja tabela relacionada NÃO foi migrada. O Count (na tabela
// principal) deve executar com sucesso, mas o Preload falhará quando o
// GORM tentar carregar os registros relacionados durante o Find.
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
