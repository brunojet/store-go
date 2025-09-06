package domain

import "testing"

func TestTerminalModel_TableName(t *testing.T) {
	var v TerminalModel
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
