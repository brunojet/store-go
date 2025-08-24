package domain

import "testing"

func TestModeloTerminal_TableName(t *testing.T) {
	var v ModeloTerminal
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
