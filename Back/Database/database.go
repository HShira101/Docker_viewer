package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Inicializar(ruta string) error {
	db, err := sql.Open("sqlite", ruta)
	if err != nil {
		return err
	}
	DB = db
	return migrar()
}

func migrar() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS contenedores_running (
			id               TEXT PRIMARY KEY,
			nombre           TEXT NOT NULL,
			imagen           TEXT NOT NULL,
			estado           TEXT NOT NULL,
			ultima_consulta  DATETIME NOT NULL
		)
	`)
	return err
}
