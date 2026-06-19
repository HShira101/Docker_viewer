package migrations

import (
	"database/sql"
	"strings"
)

// ---- Aplicar ejecuta solo las migraciones pendientes usando una tabla de tracking ----
func Aplicar(db *sql.DB) error {
	// Crea la tabla de tracking si no existe
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migraciones (
			id          INTEGER PRIMARY KEY,
			aplicada_en DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	for i, sql := range lista {
		var count int
		db.QueryRow(`SELECT COUNT(*) FROM migraciones WHERE id = ?`, i).Scan(&count)
		if count > 0 {
			continue // ya aplicada
		}

		if _, err := db.Exec(sql); err != nil {
			// SQLite no soporta ADD COLUMN IF NOT EXISTS ni RENAME IF EXISTS —
			// ignorar errores conocidos de migraciones ya aplicadas parcialmente
			msg := err.Error()
			if strings.Contains(msg, "duplicate column name") ||
				strings.Contains(msg, "already exists") ||
				strings.Contains(msg, "already another table") ||
				strings.Contains(msg, "no such table") {
				db.Exec(`INSERT INTO migraciones (id) VALUES (?)`, i)
				continue
			}
			return err
		}

		db.Exec(`INSERT INTO migraciones (id) VALUES (?)`, i)
	}
	return nil
}

var lista = []string{
	// 000 — tablas base
	`CREATE TABLE IF NOT EXISTS contenedores_running (
		id               TEXT PRIMARY KEY,
		nombre           TEXT NOT NULL,
		imagen           TEXT NOT NULL,
		estado           TEXT NOT NULL,
		ultima_consulta  DATETIME NOT NULL
	)`,
	// 001
	`CREATE TABLE IF NOT EXISTS usuarios (
		nombre   TEXT PRIMARY KEY,
		password TEXT NOT NULL
	)`,
	// 002 — agrega proyecto Compose
	`ALTER TABLE contenedores_running ADD COLUMN compose_project TEXT NOT NULL DEFAULT ''`,
	// 003 — agrega puertos publicados como JSON
	`ALTER TABLE contenedores_running ADD COLUMN puertos TEXT NOT NULL DEFAULT '[]'`,
	// 004 — renombra tabla a nombre genérico sin sufijo _running
	`ALTER TABLE contenedores_running RENAME TO contenedores`,
	// 005 — guarda el timestamp del último log enviado a VictoriaLogs
	`ALTER TABLE contenedores ADD COLUMN ultimo_log_guardado DATETIME`,

	// 006 — limpia la tabla fantasma creada por el sistema de migraciones anterior
	`DROP TABLE IF EXISTS contenedores_running`,
}
