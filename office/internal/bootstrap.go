package internal

import (
	"github.com/brunojet/store-go/office/internal/service"
	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Bootstrap creates services wired with infra repositories using a real Postgres DB.
func Bootstrap() (*service.CategoriaService, *service.TipoCategoriaService, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable search_path=store_go"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	// run minimal AutoMigrate for the types we will use
	if err := db.AutoMigrate(&domain.TipoCategoria{}, &domain.Categoria{}); err != nil {
		return nil, nil, err
	}

	// construct repositories using shared public constructors
	catRepo := repo.NewCategoriaRepo(db)
	tcRepo := repo.NewTipoCategoriaRepo(db)

	catSvc := service.NewCategoriaService(catRepo)
	tcSvc := service.NewTipoCategoriaService(tcRepo)
	return catSvc, tcSvc, nil
}
