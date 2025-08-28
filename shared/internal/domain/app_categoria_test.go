package domain

import (
	"testing"
)

func TestAppCategoria_TableName(t *testing.T) {
	var v CadastroCategoria
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
