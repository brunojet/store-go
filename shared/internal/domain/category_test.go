package domain

import "testing"

func TestCategory_TableName(t *testing.T) {
	var v Category
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
