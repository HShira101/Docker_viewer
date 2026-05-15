package controladores

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	sesion "docker_viewer/front/Sesion"
)

// ---- Consulta al backend y devuelve la lista de contenedores ----
func obtenerContenedores() []Contenedor {
	// --- Lee la URL del backend desde variable de entorno ---
	backURL := os.Getenv("BACK_URL")
	if backURL == "" {
		backURL = "http://back:10001" // ← valor por defecto dentro de Docker
	}

	resp, err := http.Get(backURL + "/api/contenedores")
	if err != nil {
		log.Println("Back API:", err)
		return nil
	}
	defer resp.Body.Close()

	var lista []Contenedor
	if err := json.NewDecoder(resp.Body).Decode(&lista); err != nil {
		log.Println("Decode contenedores:", err)
		return nil
	}
	return lista
}

// ---- Reenvía una acción POST al backend y copia la respuesta JSON al cliente ----
func (c *Controlador) proxyAccion(w http.ResponseWriter, r *http.Request, accion string) {
	// --- Recibe accion como string "iniciar" o "detener" ---
	backURL := os.Getenv("BACK_URL")
	if backURL == "" {
		backURL = "http://back:10001"
	}

	id := r.PathValue("id")
	resp, err := http.Post(backURL+"/api/contenedores/"+id+"/"+accion, "application/json", nil)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]any{"ok": false, "mensaje": "Backend no disponible"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body) // ← copia la respuesta del back directamente al browser
}

// ---- Renderiza la vista de contenedores con los datos del backend ----
func (c *Controlador) MostrarContenedores(w http.ResponseWriter, r *http.Request) {
	datos := DatosContenedores{
		DatosLayout:  DatosLayout{NombreUsuario: sesion.ObtenerUsuario(r), PaginaActual: "contenedores", CSS: "contenedores.css"},
		Contenedores: obtenerContenedores(),
	}
	c.plantillas.Contenedores.ExecuteTemplate(w, "layout.html", datos)
}

// ---- Devuelve la lista de contenedores como JSON para el fetch del cliente ----
func (c *Controlador) APIContenedores(w http.ResponseWriter, r *http.Request) {
	// --- Usado por el botón Actualizar en contenedores.html ---
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obtenerContenedores())
}

// ---- Proxy de POST /api/contenedores/{id}/iniciar → back ----
func (c *Controlador) APIIniciar(w http.ResponseWriter, r *http.Request) {
	c.proxyAccion(w, r, "iniciar")
}

// ---- Proxy de POST /api/contenedores/{id}/detener → back ----
func (c *Controlador) APIDetener(w http.ResponseWriter, r *http.Request) {
	c.proxyAccion(w, r, "detener")
}
