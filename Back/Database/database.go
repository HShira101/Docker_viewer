package database

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"

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

// ---- Crea las tablas necesarias y siembra el usuario inicial si no existe ----
func migrar() error {
	if err := crearTablas(); err != nil {
		return err
	}
	return sembrarUsuario()
}

// ---- Crea las tablas si no existen ----
func crearTablas() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS contenedores_running (
			id               TEXT PRIMARY KEY,  -- ID único del contenedor Docker
			nombre           TEXT NOT NULL,
			imagen           TEXT NOT NULL,
			estado           TEXT NOT NULL,
			ultima_consulta  DATETIME NOT NULL  -- fecha/hora de la última sincronización con Docker
		);

		CREATE TABLE IF NOT EXISTS usuarios (
			nombre     TEXT PRIMARY KEY,
			password   TEXT NOT NULL  -- sha256 de la contraseña
		);
	`)
	return err
}

// ---- Inserta el usuario desde las variables de entorno si no existe ----
func sembrarUsuario() error {
	// --- Lee credenciales desde APP_USUARIO y APP_PASSWORD ---
	nombre := os.Getenv("APP_USUARIO")
	password := os.Getenv("APP_PASSWORD")
	if nombre == "" || password == "" {
		return nil // ← sin variables de entorno no hace nada
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password))) // ← guarda el hash, nunca la contraseña plana

	_, err := DB.Exec(`
		INSERT INTO usuarios (nombre, password)
		VALUES (?, ?)
		ON CONFLICT(nombre) DO NOTHING
	`, nombre, hash)
	return err
}
