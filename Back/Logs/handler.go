package logs

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	database "docker_viewer/back/Database"
)

// ---- Handler GET /api/logs/query?id=X&limite=N ----
// Recolecta logs frescos del contenedor y consulta VictoriaLogs
func Query(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "falta parámetro id", http.StatusBadRequest)
		return
	}

	limite := 200
	if l := r.URL.Query().Get("limite"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limite = n
		}
	}

	// Busca el nombre del contenedor en SQLite
	var nombre string
	err := database.DB.QueryRow(`SELECT nombre FROM contenedores WHERE id = ?`, id).Scan(&nombre)
	if err != nil {
		http.Error(w, "contenedor no encontrado", http.StatusNotFound)
		return
	}

	// Recolecta logs frescos antes de consultar
	Recolectar(id, nombre)

	// Consulta VictoriaLogs con LogsQL
	lineas, err := consultarVictoriaLogs(id, limite)
	if err != nil {
		http.Error(w, "error consultando logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lineas)
}

// ---- Handler GET /api/logs/stream/{id} ----
// Abre una conexión SSE con logs en tiempo real del contenedor
func Stream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "falta id del contenedor", http.StatusBadRequest)
		return
	}

	var nombre string
	err := database.DB.QueryRow(`SELECT nombre FROM contenedores WHERE id = ?`, id).Scan(&nombre)
	if err != nil {
		http.Error(w, "contenedor no encontrado", http.StatusNotFound)
		return
	}

	StreamLogs(r.Context(), id, nombre, w)
}

// ---- Consulta VictoriaLogs y devuelve las últimas N líneas del contenedor ----
func consultarVictoriaLogs(id string, limite int) ([]LineaLog, error) {
	query := `{container_id="` + id + `"}`
	url := urlVlogs() + "/select/logsql/query?query=" + query + "&limit=" + strconv.Itoa(limite)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cuerpo, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// VictoriaLogs responde en JSONL (una línea JSON por log)
	var lineas []LineaLog
	for _, linea := range strings.Split(strings.TrimSpace(string(cuerpo)), "\n") {
		if linea == "" {
			continue
		}
		var ll LineaLog
		if err := json.Unmarshal([]byte(linea), &ll); err == nil {
			lineas = append(lineas, ll)
		}
	}

	return lineas, nil
}
