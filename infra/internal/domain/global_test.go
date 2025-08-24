package domain

import (
	"testing"

	"gorm.io/gorm"
)

func TestAnexo_TableName(t *testing.T) {
	a := Anexo{}
	if got := a.TableName(); got != "anexo" {
		t.Fatalf("expected table name 'anexo', got '%s'", got)
	}
}

func TestBaseEntity_BeforeCreate_Valid(t *testing.T) {
	e := &BaseEntity{Nome: "ok", Ativo: true}
	// pass nil DB (method only validates Nome length)
	if err := e.BeforeCreate(nil); err != nil {
		t.Fatalf("expected no error for valid Nome, got: %v", err)
	}
}

func TestBaseEntity_BeforeCreate_Invalid(t *testing.T) {
	e := &BaseEntity{Nome: "", Ativo: true}
	if err := e.BeforeCreate(&gorm.DB{}); err == nil {
		t.Fatalf("expected error for empty Nome, got nil")
	}
}
