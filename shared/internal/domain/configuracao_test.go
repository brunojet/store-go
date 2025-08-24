package domain

import "testing"

func TestConfiguracao_TableName(t *testing.T) {
	var v Configuracao
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
