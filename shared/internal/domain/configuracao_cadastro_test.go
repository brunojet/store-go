package domain

import "testing"

func TestConfiguracaoCadastro_TableName(t *testing.T) {
	var v ConfiguracaoCadastro
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
