package internal

import (
	"github.com/brunojet/store-go/office/internal/service"
	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/brunojet/store-go/shared/pkg/repo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Bootstrap creates services wired with infra repositories using a real Postgres DB.
func Bootstrap() (*service.CategoryService, *service.CategoryTypeService, *service.ApplicationService, error) {
	dsn := "root:p0o9i8u7@tcp(127.0.0.1:3306)/store?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, nil, err
	}

	// Executar pós-migrações para entidades que implementam PostMigratable
	if err := domain.AutoMigrate(db); err != nil {
		return nil, nil, nil, err
	}

	// construct repositories using shared public constructors
	catRepo := repo.NewCategoryRepo(db)
	tcRepo := repo.NewCategoryTypeRepo(db)
	appsRepo := repo.NewApplicationRepo(db)

	catSvc := service.NewCategoryService(catRepo)
	tcSvc := service.NewCategoryTypeService(tcRepo)
	appsSvc := service.NewApplicationService(appsRepo)
	return catSvc, tcSvc, appsSvc, nil
}
