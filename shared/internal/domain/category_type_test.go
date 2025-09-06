package domain

import "testing"

func TestCategoryType_TableName(t *testing.T) {
	var v CategoryType
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}
