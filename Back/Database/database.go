package database

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"

	migrations "docker_viewer/back/Database/Migrations"
	_ "modernc.org/sqlite"
)

// ---- Instancia global de la base de datos compartida entre paquetes ----
var DB *sql.DB

// ---- Abre la conexión con el archivo SQLite y ejecuta las migraciones ----
func Inicializar(ruta string) error {
	db, err := sql.Open("sqlite", ruta)
	if err != nil {
		return err
	}
	DB = db
	if err := migrations.Aplicar(DB); err != nil {
		return err
	}
	return sembrarUsuario()
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
