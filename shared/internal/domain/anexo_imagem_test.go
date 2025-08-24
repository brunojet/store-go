package domain_test

import (
	. "store-go/shared/internal/domain"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestDeleteAnexoCascadeImagem(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable foreign key enforcement: %v", err)
	}

	if err := db.AutoMigrate(&Anexo{}, &Imagem{}, &VersaoAplicativo{}); err != nil {
		t.Fatalf("failed to migrate tables: %v", err)
	}

	// Criar Anexo
	anexo := Anexo{
		Nome:          "Teste",
		TipoMime:      "image/png",
		MD5:           "abc123",
		Tamanho:       1234,
		Armazenamento: "s3",
		Caminho:       "imagem/1",
		Presente:      true,
	}
	if res := db.Create(&anexo); res.Error != nil {
		t.Fatalf("failed to create anexo: %v", res.Error)
	}

	imagem := Imagem{
		ID: anexo.ID,
	}
	if res := db.Create(&imagem); res.Error != nil {
		t.Fatalf("failed to create imagem: %v", res.Error)
	}

	res := db.Delete(&anexo)
	if res.Error == nil {
		t.Fatalf("esperado erro ao deletar Anexo com Imagem vinculada, mas delete teve sucesso")
	}
	t.Logf("delete de anexo falhou como esperado: %v", res.Error)

	res = db.Delete(&imagem)
	if res.Error != nil {
		t.Fatalf("failed to delete imagem: %v", res.Error)
	}

	var imgCount int64
	if r := db.Model(&Imagem{}).Where("id = ?", anexo.ID).Count(&imgCount); r.Error != nil {
		t.Fatalf("failed to count imagens: %v", r.Error)
	}
	if imgCount != 0 {
		t.Fatalf("esperado imagem existir após delete falhar, mas não foi encontrada")
	}

	if r := db.Delete(&anexo); r.Error != nil {
		t.Fatalf("failed to delete anexo after imagem removed: %v", r.Error)
	}
}
