package domain_test

import (
	"testing"

	"github.com/brunojet/store-go/shared/internal/domain"
)

func TestCatalogoAplicativo_TableName(t *testing.T) {
	var v domain.CatalogoAplicativo
	if got := v.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}
}

// // Função utilitária para atualizar id_configuracao nos registros de CatalogoAplicativo onde está nulo
// func AtualizaIdConfiguracaoCatalogoAplicativo(db *gorm.DB) error {
// 	var catalogos []domain.CatalogoAplicativo
// 	if err := db.Where("id_configuracao IS NULL").Find(&catalogos).Error; err != nil {
// 		return err
// 	}
// 	for _, cat := range catalogos {
// 		var cfg domain.Configuracao
// 		err := db.Where("id_modelo_terminal = ? AND id_tipo_integracao = ? AND id_app = ?",
// 			cat.IdModeloTerminal, cat.IdTipoIntegracao, cat.IdApp).First(&cfg).Error
// 		if err == nil {
// 			db.Model(&domain.CatalogoAplicativo{}).
// 				Where("id = ?", cat.ID).
// 				Update("id_configuracao", cfg.ID)
// 		}
// 	}
// 	return nil
// }

// func TestAtualizaIdConfiguracaoCatalogoAplicativo(t *testing.T) {
// 	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable search_path=store_go"
// 	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	err := AtualizaIdConfiguracaoCatalogoAplicativo(db)
// 	if err != nil {
// 		t.Errorf("Erro ao atualizar id_configuracao: %v", err)
// 	}
// }

// func AtualizaIdCadastroCatalogoAplicativo(db *gorm.DB) error {
// 	var catalogos []domain.CatalogoAplicativo
// 	if err := db.Where("id_cadastro IS NULL").Find(&catalogos).Error; err != nil {
// 		return err
// 	}
// 	for _, cat := range catalogos {
// 		var versao domain.VersaoAplicativo
// 		err := db.Where("id = ?", cat.IdVersaoAplicativo).First(&versao).Error
// 		if err == nil {
// 			db.Model(&domain.CatalogoAplicativo{}).
// 				Where("id = ?", cat.ID).
// 				Update("id_cadastro", versao.IdCadastro)
// 		}
// 	}
// 	return nil
// }

// func TestAtualizaIdCadastroCatalogoAplicativo(t *testing.T) {
// 	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable search_path=store_go"
// 	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	err := AtualizaIdCadastroCatalogoAplicativo(db)
// 	if err != nil {
// 		t.Errorf("Erro ao atualizar id_cadastro: %v", err)
// 	}
// }
