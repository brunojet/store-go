package repo

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupDryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open dry db: %v", err)
	}
	return db
}

func TestBuildPageable(t *testing.T) {
	tests := []struct {
		name       string
		p          *ListParams
		wantLimit  int
		wantOffset int
	}{
		{"nil", nil, DefaultLimit, 0},
		{"zero limit", &ListParams{Limit: 0, Offset: 0}, DefaultLimit, 0},
		{"small limit", &ListParams{Limit: 10, Offset: 5}, 10, 5},
		{"big limit", &ListParams{Limit: MaxLimit + 100, Offset: 2}, MaxLimit, 2},
		{"neg offset", &ListParams{Limit: 20, Offset: -5}, 20, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, o := buildPageable(tt.p)
			if l != tt.wantLimit || o != tt.wantOffset {
				t.Fatalf("got limit=%d offset=%d want limit=%d offset=%d", l, o, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}

func TestBuildOrderAndApplyListParams(t *testing.T) {
	db := setupDryDB(t)

	// valid single order via OrderColumn/OrderDir
	p1 := &ListParams{OrderColumn: "name", OrderDir: "desc", Limit: 5}
	// assert buildOrder directly
	if col, desc := buildOrder(p1); col != "name" || !desc {
		t.Fatalf("buildOrder p1 failed, got col=%q desc=%v", col, desc)
	}

	// invalid column should return empty
	p2 := &ListParams{OrderColumn: "badcol", OrderDir: "asc", Limit: 5}
	if col, _ := buildOrder(p2); col != "" {
		t.Fatalf("expected empty for invalid col")
	}

	// multiple via Order string (comma is not supported yet, but single works)
	p3 := &ListParams{Order: "name asc", Limit: 3}
	if col, desc := buildOrder(p3); col != "name" || desc {
		t.Fatalf("buildOrder p3 wrong, got col=%q desc=%v", col, desc)
	}

	// applyListParams should produce SQL with ORDER, LIMIT and OFFSET
	p4 := &ListParams{OrderColumn: "name", OrderDir: "asc", Limit: 2, Offset: 1}
	// ensure applyListParams runs and returns a non-nil DB
	q := applyListParams(db.Table("test_models"), p4)
	if q == nil {
		t.Fatalf("applyListParams returned nil")
	}
	// also call with nil params to ensure defaults apply (should not panic)
	q2 := applyListParams(db.Table("test_models"), nil)
	if q2 == nil {
		t.Fatalf("applyListParams returned nil for nil params")
	}
}

func TestBuildOrderBranches(t *testing.T) {
	// p == nil -> default
	col, _ := buildOrder(nil)
	if col != DefaultOrderColumn {
		t.Fatalf("expected default col %s got %s", DefaultOrderColumn, col)
	}

	// empty order/columns -> no order
	pEmpty := &ListParams{}
	if c, _ := buildOrder(pEmpty); c != "" {
		t.Fatalf("expected empty col for empty params, got %q", c)
	}

	// missing dir -> default asc
	pNoDir := &ListParams{OrderColumn: "name"}
	if c, d := buildOrder(pNoDir); c != "name" || d {
		t.Fatalf("expected name asc got %q desc=%v", c, d)
	}

	// invalid dir -> reject
	pBadDir := &ListParams{OrderColumn: "name", OrderDir: "bad"}
	if c, _ := buildOrder(pBadDir); c != "" {
		t.Fatalf("expected empty for invalid dir, got %q", c)
	}
}
