package domain

import "testing"

func TestApplication_TableName_And_BeforeCreate(t *testing.T) {
	var a Application
	if got := a.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
