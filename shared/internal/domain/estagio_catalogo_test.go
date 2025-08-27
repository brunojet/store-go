package domain

import "testing"

func TestEstagioCatalogo_TableName(t *testing.T) {
	var v Estagio
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
