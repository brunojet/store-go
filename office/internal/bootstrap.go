package internal

import (
	"github.com/brunojet/store-go/office/internal/service"
	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Bootstrap creates services wired with infra repositories using a real Postgres DB.
func Bootstrap() (*service.CategoryService, *service.CategoryTypeService, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable search_path=store_go"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	// Executa AutoMigrate em loop para cada entidade
	for _, entidade := range domain.EntidadesAutoMigrate {
		if err := db.AutoMigrate(entidade); err != nil {
			return nil, nil, err
		}
	}

	// Executar pós-migrações para entidades que implementam PostMigratable
	if err := domain.RunPostMigrations(db); err != nil {
		return nil, nil, err
	}

	// construct repositories using shared public constructors
	catRepo := repo.NewCategoryRepo(db)
	tcRepo := repo.NewCategoryTypeRepo(db)

	catSvc := service.NewCategoryService(catRepo)
	tcSvc := service.NewCategoryTypeService(tcRepo)
	return catSvc, tcSvc, nil
}
