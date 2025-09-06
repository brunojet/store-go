package domain

import (
	"fmt"

	"gorm.io/gorm"
)

// EnsureCompositePKAndFK garante PK e FK composta para qualquer tabela
// createTableSQL: SQL completo para criar a tabela corretamente (com tipos e constraints)
func EnsureCompositePKAndFK(db *gorm.DB, tableName string, pkCols []string, fkName, fkDef, createTableSQL string) error {
	dialect := db.Dialector.Name()
	// Ajusta fkDef conforme o banco
	switch dialect {
	case "postgres":
		// Para Postgres, fkDef deve ser um comando ALTER TABLE
		// (mantém o comportamento atual)
		// Verifica PK via information_schema
		var exists bool
		checkPK := `SELECT EXISTS (
			SELECT 1 FROM information_schema.table_constraints 
			WHERE constraint_type = 'PRIMARY KEY' 
			AND table_name = ? 
			AND table_schema = CURRENT_SCHEMA()
		)`
		if err := db.Raw(checkPK, tableName).Scan(&exists).Error; err != nil {
			return err
		}
		if !exists {
			pkColsStr := ""
			for i, col := range pkCols {
				if i > 0 {
					pkColsStr += ", "
				}
				pkColsStr += col
			}
			createPK := fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT pk_%s PRIMARY KEY (%s)`, tableName, tableName, pkColsStr)
			if err := db.Exec(createPK).Error; err != nil {
				return err
			}
		}
		// Garantir FK composta
		if fkDef != "" && fkName != "" {
			type fkInfo struct{ Name string }
			var fks []fkInfo
			checkFK := `SELECT constraint_name as name FROM information_schema.table_constraints WHERE table_name = ? AND constraint_type = 'FOREIGN KEY' AND constraint_name = ?`
			if err := db.Raw(checkFK, tableName, fkName).Scan(&fks).Error; err != nil {
				return err
			}
			if len(fks) == 0 {
				if err := db.Exec(fkDef).Error; err != nil {
					return err
				}
			}
		}
	case "sqlite":
		var tName string
		db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&tName)
		if tName == "" {
			// Tabela não existe: cria usando o SQL fornecido
			if createTableSQL == "" {
				return fmt.Errorf("recriação necessária: tabela %s não existe e createTableSQL não foi fornecido", tableName)
			}
			if err := db.Exec(createTableSQL).Error; err != nil {
				return fmt.Errorf("erro ao criar tabela %s: %w", tableName, err)
			}
			return nil
		}

		// Tabela existe, checa se PK composta está presente
		type colInfo struct{ PK int }
		var infos []colInfo
		db.Raw(fmt.Sprintf("PRAGMA table_info('%s')", tableName)).Scan(&infos)
		pkCount := 0
		for _, info := range infos {
			if info.PK > 0 {
				pkCount++
			}
		}

		// Checa FK
		fkMissing := false
		if fkDef != "" && fkName != "" {
			type fkInfo struct{ Name string }
			var fks []fkInfo
			db.Raw(fmt.Sprintf("PRAGMA foreign_key_list('%s')", tableName)).Scan(&fks)
			found := false
			for _, fk := range fks {
				if fk.Name == fkName {
					found = true
					break
				}
			}
			if !found {
				fkMissing = true
			}
		}

		// Se faltar PK composta ou FK, recria a tabela de uma vez usando a função que trata ambas
		if pkCount != len(pkCols) || fkMissing {
			// PK composta ou FK ausente: recria usando o SQL fornecido
			if err := RecreateTableForCompositePKAndFK(db, tableName, createTableSQL); err != nil {
				return err
			}
		}
		// findIndexInsensitive retorna o índice da substring (case-insensitive), ou -1 se não encontrar
	}
	return nil
}

// RecreateTableForCompositePKAndFK recria uma tabela SQLite usando o SQL fornecido, movendo os dados das colunas em comum
func RecreateTableForCompositePKAndFK(db *gorm.DB, tableName string, createTableSQL string) error {
	oldTable := tableName + "_old"
	// Renomeia a tabela antiga
	if err := db.Exec("ALTER TABLE " + tableName + " RENAME TO " + oldTable).Error; err != nil {
		return fmt.Errorf("erro ao renomear tabela antiga: %w", err)
	}
	// Cria a nova tabela
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("erro ao criar nova tabela %s: %w", tableName, err)
	}
	// Descobre colunas em comum
	var oldCols, newCols []string
	rows, err := db.Raw("PRAGMA table_info('" + oldTable + "')").Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cid int
			var name, ctype string
			var notnull, pk int
			var dflt interface{}
			_ = rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
			oldCols = append(oldCols, name)
		}
	}
	rows2, err2 := db.Raw("PRAGMA table_info('" + tableName + "')").Rows()
	if err2 == nil {
		defer rows2.Close()
		for rows2.Next() {
			var cid int
			var name, ctype string
			var notnull, pk int
			var dflt interface{}
			_ = rows2.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
			newCols = append(newCols, name)
		}
	}
	// Interseção
	var commonCols []string
	for _, c := range oldCols {
		for _, n := range newCols {
			if c == n {
				commonCols = append(commonCols, c)
			}
		}
	}
	if len(commonCols) > 0 {
		colsList := ""
		for i, c := range commonCols {
			if i > 0 {
				colsList += ", "
			}
			colsList += c
		}
		copyData := "INSERT INTO " + tableName + " (" + colsList + ") SELECT " + colsList + " FROM " + oldTable
		if err := db.Exec(copyData).Error; err != nil {
			return fmt.Errorf("erro ao copiar dados: %w", err)
		}
	}
	// Remove a tabela antiga
	if err := db.Exec("DROP TABLE " + oldTable).Error; err != nil {
		return fmt.Errorf("erro ao remover tabela antiga: %w", err)
	}
	return nil
}
