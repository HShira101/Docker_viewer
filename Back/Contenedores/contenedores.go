package contenedores

import (
	"encoding/json"
	"net/http"

	database "docker_viewer/back/Database"
)

// ---- Handler de GET /api/contenedores: sincroniza Docker, lee SQLite y devuelve JSON ----
func Listar(w http.ResponseWriter, r *http.Request) {
	Actualizar() // ← refresca los datos desde Docker antes de responder

	rows, err := database.DB.Query(`
		SELECT id, nombre, imagen, estado, ultima_consulta
		FROM contenedores_running
		ORDER BY nombre
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	lista := make([]Contenedor, 0)
	for rows.Next() {
		var c Contenedor
		rows.Scan(&c.ID, &c.Nombre, &c.Imagen, &c.Estado, &c.UltimaConsulta)
		lista = append(lista, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}
