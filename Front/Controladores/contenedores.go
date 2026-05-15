package controladores

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

// ---- Renderiza la vista de contenedores con los datos del backend ----
func (c *Controlador) MostrarContenedores(w http.ResponseWriter, r *http.Request) {
	datos := DatosContenedores{
		DatosLayout:  DatosLayout{NombreUsuario: "Shira", PaginaActual: "contenedores", CSS: "contenedores.css"},
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
