package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCreateImage(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable search_path=store_go"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(t, err)

	// Migrar as tabelas necessárias com FK explícita
	require.NoError(t, err)

	// Criar o objeto base
	storage := StorageObject{}
	err = db.Create(&storage).Error
	require.NoError(t, err)
	require.NotZero(t, storage.ID)

	// Criar a imagem vinculada ao objeto base
	img := Image{
		BaseModel: BaseModel{ID: storage.ID},
		ImageType: ImageTypeIcon,
	}
	err = db.Create(&img).Error
	require.NoError(t, err)

	// Buscar e validar
	var found Image
	err = db.Preload("StorageObject").First(&found, img.ID).Error
	require.NoError(t, err)
	require.Equal(t, img.ID, found.ID)
	require.Equal(t, ImageTypeIcon, found.ImageType)
	require.Equal(t, storage.ID, found.StorageObject.ID)
}
