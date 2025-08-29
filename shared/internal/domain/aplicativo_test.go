package domain

import "testing"

func TestAplicativo_TableName_And_BeforeCreate(t *testing.T) {
	var a Aplicativo
	if got := a.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
