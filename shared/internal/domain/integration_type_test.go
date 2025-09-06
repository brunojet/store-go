package domain

import "testing"

func TestIntegrationType_TableName(t *testing.T) {
	var v IntegrationType
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
