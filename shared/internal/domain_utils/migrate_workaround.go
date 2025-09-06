package domain_tools

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func CreatePrimaryKey(db *gorm.DB, tableName, constraintName string, pkFields []string) error {
	dialect := db.Dialector.Name()
	var exists bool
	var checkSQL string
	switch dialect {
	case "postgres":
		checkSQL = `SELECT EXISTS (
			SELECT 1 FROM information_schema.table_constraints
			WHERE constraint_type = 'PRIMARY KEY'
			AND table_name = ?
		);`
		if err := db.Raw(checkSQL, tableName).Scan(&exists).Error; err != nil {
			return fmt.Errorf("erro ao checar existência da PK: %w", err)
		}
	case "mysql":
		checkSQL = `SELECT COUNT(*) > 0 FROM information_schema.TABLE_CONSTRAINTS
			WHERE CONSTRAINT_TYPE = 'PRIMARY KEY'
			AND TABLE_NAME = ? AND TABLE_SCHEMA = (SELECT DATABASE())`
		if err := db.Raw(checkSQL, tableName).Scan(&exists).Error; err != nil {
			return fmt.Errorf("erro ao checar existência da PK: %w", err)
		}
	default:
		return fmt.Errorf("dialeto não suportado para verificação de PK: %s", dialect)
	}
	if exists {
		return nil
	}
	pkFieldsStr := strings.Join(pkFields, ", ")
	sql := fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s PRIMARY KEY (%s)`, tableName, constraintName, pkFieldsStr)
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("erro ao criar PK: %w", err)
	}
	return nil
}

func CreateConstraints(db *gorm.DB, tableName, tableNameReference, constraintName string, fkFields []string) error {
	dialect := db.Dialector.Name()
	var exists bool
	var checkSQL string
	switch dialect {
	case "postgres":
		checkSQL = `SELECT EXISTS (
		  SELECT 1 FROM information_schema.table_constraints
		  WHERE constraint_name = ? AND table_name = ?
		);`
		if err := db.Raw(checkSQL, constraintName, tableName).Scan(&exists).Error; err != nil {
			return fmt.Errorf("erro ao checar existência da constraint: %w", err)
		}
	case "mysql":
		var dbName string
		if err := db.Raw("SELECT DATABASE()").Scan(&dbName).Error; err != nil {
			return fmt.Errorf("erro ao obter nome do banco: %w", err)
		}
		checkSQL = `SELECT EXISTS (
		  SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
		  WHERE CONSTRAINT_NAME = ? AND TABLE_NAME = ? AND CONSTRAINT_SCHEMA = ?
		);`
		if err := db.Raw(checkSQL, constraintName, tableName, dbName).Scan(&exists).Error; err != nil {
			return fmt.Errorf("erro ao checar existência da constraint: %w", err)
		}
	default:
		return fmt.Errorf("dialeto não suportado para verificação de constraint: %s", dialect)
	}
	if exists {
		return nil
	}

	fkFieldsStr := strings.Join(fkFields, ", ")

	sql := fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s
		FOREIGN KEY (%s)
		REFERENCES %s(%s)
		ON DELETE RESTRICT ON UPDATE RESTRICT`, tableName, constraintName, fkFieldsStr, tableNameReference, fkFieldsStr)
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("erro ao criar FK composta: %w", err)
	}
	return nil
}

func RecreateTable(db *gorm.DB, tableName, createTable string, columns []string) error {
	dialect := db.Dialector.Name()
	if dialect != "sqlite" {
		return fmt.Errorf("dialeto não suportado para recriar a tabela application_version")
	}

	tableOld := tableName + "_old"
	tableNew := tableName

	columnsStr := strings.Join(columns, ", ")

	steps := []string{
		fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tableNew, tableOld),
		createTable + ";",
		fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s;", tableNew, columnsStr, columnsStr, tableOld),
		fmt.Sprintf("DROP TABLE %s;", tableOld),
	}

	for _, sql := range steps {
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("erro ao executar passo da recriação da tabela: %w (sql: %s)", err, sql)
		}
	}
	return nil
}
