package contenedores

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	database "docker_viewer/back/Database"
)

// ---- Handler de GET /api/contenedores: sincroniza Docker, lee SQLite y devuelve JSON ----
func Listar(w http.ResponseWriter, r *http.Request) {
	Actualizar() // ← refresca los datos desde Docker antes de responder

	rows, err := database.DB.Query(`
		SELECT id, nombre, imagen, estado, ultima_consulta, compose_project, puertos
		FROM contenedores_running
		ORDER BY compose_project, nombre
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	lista := make([]Contenedor, 0)
	for rows.Next() {
		var c Contenedor
		var puertosJSON string
		rows.Scan(&c.ID, &c.Nombre, &c.Imagen, &c.Estado, &c.UltimaConsulta, &c.ComposeProject, &puertosJSON)
		json.Unmarshal([]byte(puertosJSON), &c.Puertos)
		lista = append(lista, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

// ---- Handler de POST /api/contenedores/{id}/iniciar ----
func Iniciar(w http.ResponseWriter, r *http.Request) {
	// --- Recibe id del contenedor desde la URL ---
	ejecutarAccion(w, r.PathValue("id"), "start", "iniciado")
}

// ---- Handler de POST /api/contenedores/{id}/detener ----
func Detener(w http.ResponseWriter, r *http.Request) {
	// --- Recibe id del contenedor desde la URL ---
	ejecutarAccion(w, r.PathValue("id"), "stop", "detenido")
}

// ---- Comprueba existencia en DB, llama a Docker y responde con éxito o error ----
func ejecutarAccion(w http.ResponseWriter, id, accionDocker, textoExito string) {
	w.Header().Set("Content-Type", "application/json")

	// ---- Comprueba que el contenedor existe en la base de datos ----
	var nombre string
	err := database.DB.QueryRow(
		`SELECT nombre FROM contenedores_running WHERE id = ?`, id,
	).Scan(&nombre)

	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(Respuesta{OK: false, Mensaje: "Contenedor no encontrado en la base de datos"})
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(Respuesta{OK: false, Mensaje: "Error al consultar la base de datos"})
		return
	}

	// ---- Envía la acción al socket de Docker ----
	status, err := dockerPost("/containers/" + id + "/" + accionDocker)
	if err != nil {
		json.NewEncoder(w).Encode(Respuesta{OK: false, Mensaje: "Error al contactar Docker: " + err.Error()})
		return
	}

	switch status {
	case 204:
		json.NewEncoder(w).Encode(Respuesta{OK: true, Mensaje: fmt.Sprintf("Contenedor %s %s correctamente", nombre, textoExito)})
	case 304:
		json.NewEncoder(w).Encode(Respuesta{OK: false, Mensaje: fmt.Sprintf("El contenedor %s ya está en ese estado", nombre)})
	case 404:
		json.NewEncoder(w).Encode(Respuesta{OK: false, Mensaje: "Contenedor no encontrado en Docker"})
	default:
		json.NewEncoder(w).Encode(Respuesta{OK: false, Mensaje: fmt.Sprintf("Docker respondió con código %d", status)})
	}
}
