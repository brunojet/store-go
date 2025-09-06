package domain

import (
	"testing"

	"gorm.io/gorm"
)

func TestVersaoApplication_TableName_And_BeforeCreate(t *testing.T) {
	var v ApplicationVersion
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}

	// ensure BeforeCreate from embedded BaseEntity runs without error when Nome is valid
	v.BaseEntity = BaseEntity{Nome: "Valid"}
	if err := v.BeforeCreate((*gorm.DB)(nil)); err != nil {
		t.Fatalf("BeforeCreate failed: %v", err)
	}
}
