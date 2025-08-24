package domain

import "testing"

func TestTipoIntegracao_TableName(t *testing.T) {
	var v TipoIntegracao
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
