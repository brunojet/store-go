package domain

import "testing"

func TestTipoCategoria_TableName(t *testing.T) {
	var v TipoCategoria
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
