package internal

import (
	"database/sql"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"store-go/shared/internal/domain"
)

func InitDB() (*gorm.DB, error) {
	// Conexão com PostgreSQL
	dsn := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	schemaName := "store_go"

	// Cria o schema se não existir
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	_, err = sqlDB.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName)
	if err != nil {
		return nil, err
	}

	// Inicializa o GORM usando o schema
	db, err := gorm.Open(postgres.Open(dsn+"&search_path="+schemaName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
			NoLowerCase:   false,
		},
	})
	if err != nil {
		return nil, err
	}

	// Migrar todas as entidades
	db.AutoMigrate(
		&domain.ModeloTerminal{},
		&domain.TipoIntegracao{},
		&domain.TipoCategoria{},
		&domain.Categoria{},
		&domain.Aplicativo{},
		&domain.AppCategoria{},
		&domain.Configuracao{},
		&domain.Cadastro{},
		&domain.ConfiguracaoCadastro{},
		&domain.VersaoAplicativo{},
		&domain.EstagioCatalogo{},
		&domain.CatalogoAplicativo{},
	)

	return db, nil
}
