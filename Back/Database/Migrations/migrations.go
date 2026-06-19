package migrations

import (
	"database/sql"
	"strings"
)

// ---- Ejecuta todas las migraciones en orden ----
func Aplicar(db *sql.DB) error {
	for _, m := range lista {
		if _, err := db.Exec(m); err != nil {
			// SQLite no soporta ADD COLUMN IF NOT EXISTS — ignorar si la columna ya existe
			if strings.Contains(err.Error(), "duplicate column name") ||
			strings.Contains(err.Error(), "already exists") {
			continue
		}
			return err
		}
	}
	return nil
}

var lista = []string{
	// 001 — tablas base
	`CREATE TABLE IF NOT EXISTS contenedores_running (
		id               TEXT PRIMARY KEY,
		nombre           TEXT NOT NULL,
		imagen           TEXT NOT NULL,
		estado           TEXT NOT NULL,
		ultima_consulta  DATETIME NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS usuarios (
		nombre   TEXT PRIMARY KEY,
		password TEXT NOT NULL
	)`,

	// 002 — agrega proyecto Compose (ignorar error si la columna ya existe)
	`ALTER TABLE contenedores_running ADD COLUMN compose_project TEXT NOT NULL DEFAULT ''`,

	// 003 — agrega puertos publicados como JSON (ignorar error si la columna ya existe)
	`ALTER TABLE contenedores_running ADD COLUMN puertos TEXT NOT NULL DEFAULT '[]'`,

	// 004 — renombra tabla a nombre genérico sin sufijo _running
	`ALTER TABLE contenedores_running RENAME TO contenedores`,

	// 005 — guarda el timestamp del último log enviado a VictoriaLogs
	`ALTER TABLE contenedores ADD COLUMN ultimo_log_guardado DATETIME`,
}
