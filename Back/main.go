package main

import (
	"log"
	"net/http"
	"os"

	contenedores "docker_viewer/back/Contenedores"
	database "docker_viewer/back/Database"
)

func main() {
	rutaDB := os.Getenv("DB_PATH")
	if rutaDB == "" {
		rutaDB = "/data/database.sqlite"
	}

	if err := database.Inicializar(rutaDB); err != nil {
		log.Fatal("DB:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /api/contenedores", contenedores.Listar)

	log.Println("Backend en :10001 → http://localhost:10001")
	log.Fatal(http.ListenAndServe(":10001", mux))
}
