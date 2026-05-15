package main

import (
	"log"
	"net/http"
	"os"

	contenedores "docker_viewer/back/Contenedores"
	database "docker_viewer/back/Database"
)

// ---- Punto de entrada del servidor backend ----
func main() {
	// --- Lee la ruta del archivo SQLite desde variable de entorno ---
	rutaDB := os.Getenv("DB_PATH")
	if rutaDB == "" {
		rutaDB = "/data/database.sqlite" // ← valor por defecto dentro de Docker
	}

	// ---- Inicializa la base de datos y ejecuta migraciones ----
	if err := database.Inicializar(rutaDB); err != nil {
		log.Fatal("DB:", err)
	}

	// ---- Registra rutas HTTP ----
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok")) // ← usado por el healthcheck del docker-compose
	})
	mux.HandleFunc("GET /api/contenedores", contenedores.Listar)

	log.Println("Backend en :10001 → http://localhost:10001")
	log.Fatal(http.ListenAndServe(":10001", mux))
}
