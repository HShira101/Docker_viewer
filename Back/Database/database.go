package database

import (
	"database/sql"

	_ "modernc.org/sqlite" // ← driver SQLite puro Go, sin CGo
)

// ---- Instancia global de la base de datos compartida entre paquetes ----
var DB *sql.DB

// ---- Abre la conexión con el archivo SQLite y ejecuta las migraciones ----
func Inicializar(ruta string) error {
	// --- Recibe ruta como string con la ruta al archivo .sqlite ---
	db, err := sql.Open("sqlite", ruta)
	if err != nil {
		return err
	}
	DB = db
	return migrar()
}

// ---- Crea las tablas necesarias si no existen ----
func migrar() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS contenedores_running (
			id               TEXT PRIMARY KEY,  -- ID único del contenedor Docker
			nombre           TEXT NOT NULL,
			imagen           TEXT NOT NULL,
			estado           TEXT NOT NULL,
			ultima_consulta  DATETIME NOT NULL  -- fecha/hora de la última sincronización con Docker
		)
	`)
	return err
}
