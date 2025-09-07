package domain_tools

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// EnsureCascadeOnDelete ajusta a constraint de FK para ON DELETE CASCADE se não estiver assim.
// table: tabela de origem
// fkColumn: coluna da FK na origem
// refTable: tabela de destino
// refColumn: coluna referenciada
func EnsureCascadeOnDelete(db *gorm.DB, table, fkColumns, refTable, refColumns string) error {
	dialect := db.Dialector.Name()
	var fkName, onDelete string
	switch dialect {
	case "postgres":
		type fkInfo struct {
			Name            string `gorm:"column:conname"`
			DeleteRule      string `gorm:"column:confdeltype"`
			ReferencedTable string `gorm:"column:referenced_table"`
		}
		norm := func(s string) string {
			return strings.ToLower(strings.Trim(s, `"`))
		}
		tableNorm := norm(table)
		refTableNorm := norm(refTable)
		query := `SELECT c.conname, c.confdeltype, t2.relname as referenced_table
		FROM pg_constraint c
		JOIN pg_class t ON c.conrelid = t.oid
		JOIN pg_class t2 ON c.confrelid = t2.oid
		WHERE lower(t.relname) = ? AND lower(t2.relname) = ? AND c.contype = 'f' LIMIT 1`
		var fk fkInfo
		if err := db.Raw(query, tableNorm, refTableNorm).Scan(&fk).Error; err != nil {
			return fmt.Errorf("erro ao buscar FKs: %w", err)
		}
		if fk.ReferencedTable != refTable || fk.DeleteRule == "c" {
			return nil
		}
		fkName = fk.Name
	case "mysql":
		query := `SELECT rc.CONSTRAINT_NAME, rc.DELETE_RULE
		FROM information_schema.REFERENTIAL_CONSTRAINTS rc
		WHERE rc.TABLE_NAME = ? AND rc.REFERENCED_TABLE_NAME = ? AND rc.CONSTRAINT_SCHEMA = (SELECT DATABASE()) LIMIT 1;`
		if err := db.Raw(query, table, refTable).Row().Scan(&fkName, &onDelete); err != nil {
			return fmt.Errorf("erro ao buscar nome da constraint e delete_rule: %w", err)
		}
	default:
		return fmt.Errorf("dialeto não suportado para alteração de FK: %s", dialect)
	}

	if err := dropFk(db, table, fkName); err != nil {
		return err
	}

	if err := createFk(db, table, fkName, fkColumns, refTable, refColumns); err != nil {
		return err
	}

	return nil
}

func dropFk(db *gorm.DB, table, fkName string) error {
	dialect := db.Dialector.Name()
	var drop string
	switch dialect {
	case "postgres":
		drop = fmt.Sprintf(`ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s`, table, fkName)
	case "mysql":
		drop = fmt.Sprintf(`ALTER TABLE %s DROP FOREIGN KEY %s`, table, fkName)
	}
	if err := db.Exec(drop).Error; err != nil {
		return fmt.Errorf("erro ao remover constraint antiga: %w", err)
	}
	return nil
}

func createFk(db *gorm.DB, table, fkName, fkColumns, refTable, refColumns string) error {
	add := fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE CASCADE ON UPDATE RESTRICT`, table, fkName, fkColumns, refTable, refColumns)
	if err := db.Exec(add).Error; err != nil {
		return fmt.Errorf("erro ao criar constraint com CASCADE: %w", err)
	}
	return nil
}
