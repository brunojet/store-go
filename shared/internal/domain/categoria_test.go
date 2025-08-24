package domain

import "testing"

func TestCategoria_TableName(t *testing.T) {
	var v Categoria
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
