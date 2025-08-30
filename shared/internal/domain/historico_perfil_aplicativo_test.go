package domain

import "testing"

func TestCadastro_TableName(t *testing.T) {
	var v HistoricoPerfilAplicativo
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
