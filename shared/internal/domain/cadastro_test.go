package domain

import "testing"

func TestCadastro_TableName(t *testing.T) {
	var v Cadastro
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
