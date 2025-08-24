package domain

import "testing"

func TestCatalogoAplicativo_TableName(t *testing.T) {
	var v CatalogoAplicativo
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
